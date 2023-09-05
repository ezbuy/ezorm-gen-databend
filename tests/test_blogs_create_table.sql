USE `test`;

CREATE TABLE `blogs` (
	id BIGINT NOT NULL ,
	user_id INT NOT NULL ,
	title VARCHAR NOT NULL ,
	content VARCHAR NOT NULL ,
	status INT NOT NULL ,
	readed INT NOT NULL ,
	created_at TIMESTAMP NOT NULL ,
	updated_at TIMESTAMP NOT NULL 
)