### Sample ruby on rails application

```
$ cd /path/to/repos
$ cd sample_app_reference
$ bundle install --without production
```

Next, migrate the database:

```
$ bundle exec rails db:create
$ bundle exec rails db:migrate
```

Now, you'll be ready to run the app in a local server:

```
$ bundle exec rails server
```

Add the gem `appoptics_apm` to your Gemfile

```
$ bundle exec rails generate appoptics_apm:install
```
