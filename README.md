# if_inventory

### Database tables

spacecraft
```
id | name | class | crew | deleted | created
```

armourments
```
id | title | spacecraft_id | qty
```

Armourments associates to the spacecraft table by spacecraft_id.

### Structure

#### Models
 
These are classes that directly interact with the database, all database logic should be in here

#### Services

This is the http service for serving data/processing updates and insert requests
