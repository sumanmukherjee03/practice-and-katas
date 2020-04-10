DOCKER_CONTAINER := testmysql

SAMPLE_DB_GITHUB_REPO := git@github.com:datacharmer/test_db.git

docker-start:
	docker run --name testmysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=welcome2mysql -d mysql:5.7

docker-stop:
	docker stop testmysql
	docker rm testmysql

docker-clean:
	docker stop testmysql
	sleep 5
	docker rm testmysql

clean:
	rm -rf tmp

load: clean
	mkdir -p tmp \
		&& cd tmp \
		&& git clone $(SAMPLE_DB_GITHUB_REPO) \
		&& cd test_db \
		&& mysql -h localhost --user=root --password=welcome2mysql < employees.sql
	rm -rf tmp
	echo "DONE"
