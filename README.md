## Introduction
The loan application system..

## Walkthrough
This is a 3-tier application consisting of:

- A VueJs/Vuetify frontend. The implementation is under the [web] folder.
- A Golang API. The implementation is under the [api] folder.
- A Postgres database. The relevant folder is [db] folder.

### Frontend
Once a user navigates to the frontend (by default found at http://localhost:8080), if he is not logged in, he will be redirected to a login page where he will have to provide his username and password. There are a few test users registered already:

| **Username** | **Password** |
|--------------|--------------|
| Alice        | password123  |
| Bob          | password123  |
| Charlie      | password123  |
| David        | password123  |

Once logged in, he will be redirected to the loan calculator page, where he can input his business details in order to get the balance sheet. Later, after reviewing his balance sheet he can check the loan permissions.

### API
The API exposes 4 endpoints : login, providers, balance-sheet, calculate-loan

Right now, balance-sheet is fetched from a json file as accounting provider's details were not given. (For now, I commented the code for fetching balance sheet from accounting provider)

There are a few dummy balance sheets saved which is mapped by business name. Please use these business name while generating blance sheet: 
| **Business Name** |**Business Name**  |
|-------------------|-------------------|
| business1         | [business1.json](api/balanceSheet/business1.json)    |
| business2         | business1.json    |
| business3         | business1.json    | 

### Steps to run the app

    docker-compose up
Spins up the application

