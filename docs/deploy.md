# Deployment

[Heroku](https://dashboard.heroku.com/apps) handles github push webhook, recognize runtime platform and redeploy project.
Frontend and daemon divided to 2 repo's to provide defined behavior for Heroku deployments.

- [Frontend](https://frogdb-frontend.herokuapp.com/dashboard)
- [Backend](https://frogdb.herokuapp.com/docs)

Preview backend:
![Swagger page](../img/swagger.png)

# Frontend

Just `React` with code-generated REST client. Tool for generation `openapitools/openapi-generator-cli`

Regenerate client:

```bash
make frontend-sdk-gen
```

Preview frontend:

- Dashboard
  ![Dashboard](../img/f2.png)
- Create table
  ![Create table](../img/f1.png)
- Select
  ![Select](../img/f3.png)
