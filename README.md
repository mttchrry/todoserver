# TODO - A REST API for organizing tasks

This repository will serve as implementation and documentation for the rest api web backend powering the TODO app.


## Documentation

To get this system up and running you will need postgresql installed and do a pg_restore of the included schema.sql file.  You'll also need to copy the sample.env file into a new file named .env, replacing the variables with the name and location of DB you created, as well as your desired Jwt_token secret.  From there, just 
```bash 
$ go get
```
```bash
$ go build
``` 

from this directory and then run 
```bash
$ ./todo
```

and your server should be up and running, ready to serve port http://localhost:8000

## Design Questions and Concers

There are a few ambiguities in the requirements for this, and as such I want to state my assumptions that would need verified.  They are:
* User Sign in should be part of the workflow, to authenticate task creation
* There will need to be a hosted datastore for this (my implementation is a PostgreSQL db.)
* Tasks will not be removed/deleted if their owning users are deleted. 
* Tasks would likely need a status field in the future as well. 
* The ability to delete users would be limited to some power user group, not implemented here for the sake of brevity. 


## Further Improvements to consider

* allow only creators of tasks/projects update or delete them
* Add state on tasks (Not Started, In Progress, Completed)
* Add authority role, and without it don't allow update of users outside of a registered user updating themselves. 
* need to check expiration of JWT for validation, shorten the time too, or move to OAuth or the like
* allow filtering on multiple items, by names, or all tasks with due dates before a give date

## Testing
I skimped out some of the testing out of the sake of time and architecture (mocking the postgresql DB, handling RESTApi's more cleanly, etc) and that this code won't be maintained.  ideally you'd have a testing framework and at least more tests around the protected endpoints after signing in. 
