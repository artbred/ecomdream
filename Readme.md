## This project is a simple REST API that allows to train dreambooth model, charge for usage and inference it. It works on top of [Replicate](https://replicate.com)

### This project has 3 services
+ [api](./src/services/api) - rest api
+ [cron](./src/services/cron) - performs simple cron tasks
+ [imager](./src/services/imager) - grpc service for image processing

### I used following technologies
+ Redis
+ Postgresql
+ Flyway (used for database migrations)
+ gRpc
+ Prometheus and Grafana
+ Docker

#### I will describe user experience

1. User sends POST request to ``/api/v1/payments/create``, with following json body
   1. ``plan_id`` - Required. ID for plan, they can vary in amount of images that user can generate and other possible features
   2. ``promocode_id`` - Optional. Promocode ID in Stripe, if it is incorrect, payment link is created without promocode
   3. ``version_id`` - Required if ``plan.is_init is true``. All plans have column [is_init](./sql/V2__plans.sql) which is used to identify if user purchases an add-on, for example extra amount of images to generate. If this values is not set or is incorrect user gets and error, otherwise webhook receives payment confirmation and plan with extra features binds to provided ``version_id``

2. User pays using payment link and gets ``payment_id``, then one prepares images and sends request to ``/api/v1/versions/train/{id}`` where ``{id}`` is ``payment_id`` that one received after payment
   1. I have [middleware](./src/services/api/internal/middleware/freeze.go) that sets key in redis in order to block extra requests from same user (payment_id). It takes some time to prepare data and receive success response from replicate (meaning that they started the training process) so it is possible to abuse the system. After endpoint is done key will be deleted from redis, or it will be automatically deleted in 5 minutes
   2. Images are send to [imager](./src/services/imager) where they are check for being more than 512x512 pixels and that they have supported content-type, if one image does not meet requirements data about that is stored in response so user might receive all issues after first request. The process of checking involves [concurrency](https://github.com/artbred/ecomdream/blob/a9384e29da19f5a75808b11427f613865b23b7a6/src/services/imager/server.go#L17-L49)
   3. If everything is alright imager forms zip archive and uploads it to bucket, after request is send to replicate and user gets ``version_id``

3. With ``version_id`` user can call ``/api/v1/versions/info/{id}`` and get info about version status. If version is ready user gets [extended](https://github.com/artbred/ecomdream/blob/a9384e29da19f5a75808b11427f613865b23b7a6/src/domain/models/versions.go#L150-L185) info about version. At the same time [cron](./src/services/cron) is running [task](https://github.com/artbred/ecomdream/blob/a9384e29da19f5a75808b11427f613865b23b7a6/src/services/cron/jobs/push_versions/logic.go#L10-L36) which gets running version from postgres and checks to see if it is ready.

4. Once version has field ``pushed_at is not null`` user can perform request to trained model and inference prompts using endpoint ``/api/v1/prompts/create/{id}``
    1. When user sends request this endpoint also freezes using same [middleware](./src/services/api/internal/middleware/freeze.go) so there can not be concurrent request for same model
   2. When you inference any model on replicate it returns you ``prediction_id`` which is then used for getting info about running prediction. In order to continuously check for prediction status I use [this function](https://github.com/artbred/ecomdream/blob/a9384e29da19f5a75808b11427f613865b23b7a6/src/domain/replicate/predictions.go#L75-L121)
   3. After the prediction is done I transfer images from [replicate to cloudflare](https://github.com/artbred/ecomdream/blob/aae457464b3f2f0f38d62e97b20ddc1561df5c58/src/services/api/core/v1/prompts/images.go#L12-L38) using similar solution that you have seen before, however I added ``sync.WaitGroup`` since I do not need to process images until each image is uploaded to cloudflare.

5. User can send request to ``/api/v1/prompts/list/{id}`` and get info about completed prompts, in order to achieve that in one query I implemented [custom type](https://github.com/artbred/ecomdream/blob/a9384e29da19f5a75808b11427f613865b23b7a6/src/domain/models/images.go#L25-L51)

#### There are no tests in this code, not because I ignore them, but because it was supposed to be an MVP and I wanted to write the code as quickly as possible

#### The big problem with the number of users would be the replicate rate limit, to solve this problem I would make a separate service that would take requests from others and use [leaky bucket algorithm](https://github.com/uber-go/ratelimit)
