## This project is a simple REST API that allows to train dreambooth model, charge for usage and inference it. It works on top of [Replicate](https://replicate.com)

**Please note that this project is not perfect and has some drawbacks, however I believe that the current state will be enough to assess my skills. Some inefficiencies are marked in TODOs**

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

#### I will describe user experience in order to save your time

1. User sends POST request to ``/api/v1/payments/create``, with following json body
   1. ``plan_id`` - Required. ID for plan, they can vary in amount of images that user can generate and other possible features
   2. ``promocode_id`` - Optional. Promocode ID in Stripe, if it is incorrect, payment link is created without promocode
   3. ``version_id`` - Required if ``plan.is_init is true``. All plans have column [is_init](./sql/V2__plans.sql) which is used to identify if user purchases an add-onn, for example extra amount of images to generate. If this values is not set or is incorrect  user gets and error, otherwise webhook receives payment confirmation and plan with extra features binds to provided ``version_id``

2. User pays using payment link and gets ``payment_id``, then one prepares images and sends request to ``/api/v1/versions/train/{id}`` where ``{id}`` is ``payment_id`` that one received after payment
   1. I have [middleware](./src/services/api/internal/middleware/freeze.go) that sets key in redis in order to block extra requests from same user (payment_id). It takes some time to prepare data and receive success response from replicate (meaning that they started the training process) so it is possible to abuse the system. After endpoint is done key will be deleted from redis, or it will be automatically deleted in 5 minutes
   2. Images are send to [imager](./src/services/imager) where they are check for being more than 512x512 pixels and that they have supported content-type, if one image does not meet requirements it is stored in response so user might receive all issues after first request. The process of checking
