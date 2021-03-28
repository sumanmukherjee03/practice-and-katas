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

Now, you'll be ready to run the app in a local server either in interactive or in daemon mode
```
$ bundle exec rails server
$ bundle exec rails server -d
```

To stop the daemonized rails server
```
kill -9 $(cat tmp/pids/server.pid)
```

Add the gem `appoptics_apm` to your Gemfile

```
$ bundle exec rails generate appoptics_apm:install
```

```
vim config/initializers/appoptics_apm.rb
AppOpticsAPM::Config[:service_key] = '<token>:rails_sample_app'
```
