DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS userlike;
DROP TABLE IF EXISTS userlike_comment;
DROP TABLE IF EXISTS userrole;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS post_category;
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS user_activity;
DROP TABLE IF EXISTS reports;

CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		name TEXT UNIQUE,
		password TEXT,
		email TEXT UNIQUE,
		permissions TEXT,
		role_id INTEGER,
    	CONSTRAINT fk_role FOREIGN KEY(role_id) REFERENCES userrole(id) ON DELETE CASCADE
	  );

CREATE TABLE IF NOT EXISTS userrole (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);

CREATE TABLE IF NOT EXISTS category (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		description TEXT
	  );


CREATE TABLE IF NOT EXISTS post (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		title TEXT,
		body TEXT,
		tags TEXT,
		image TEXT,
		user_id INTEGER,
		category_id INTEGER,
		posting_time TEXT,
    	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    	CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE
		
	  );

CREATE TABLE IF NOT EXISTS userlike (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		mark BOOLEAN, 
		user_id INTEGER,
		post_id INTEGER,
    	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    	CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
		CONSTRAINT unique_like UNIQUE(user_id,post_id)
	);

CREATE TABLE IF NOT EXISTS userlike_comment (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		mark BOOLEAN, 
		user_id INTEGER,
		comment_id INTEGER,
    	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
    	CONSTRAINT fk_comment FOREIGN KEY(comment_id) REFERENCES comment(id) ON DELETE CASCADE,
		CONSTRAINT unique_like UNIQUE(user_id,comment_id)
	);
       
CREATE TABLE IF NOT EXISTS comment (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		post_id INTEGER,
		posting_time TEXT,
    	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE,
		CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE
	);

CREATE TABLE IF NOT EXISTS tag (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		tagname TEXT UNIQUE,
		post_id INTEGER,
		CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE
	);

CREATE TABLE IF NOT EXISTS post_tag(
    post_id INTEGER,
    tag_id INTEGER,
    CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
    CONSTRAINT fk_tag FOREIGN KEY(tag_id) REFERENCES tag(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS session(
		cookiename TEXT NOT NULL,
		cookievalue TEXT NOT NULL,
		user_id INTEGER,
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
	
);

CREATE TABLE IF NOT EXISTS reports(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		status TEXT,
		admin_reply TEXT
);

CREATE TABLE IF NOT EXISTS requests(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER
);

CREATE TABLE IF NOT EXISTS user_activity(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		notification_type TEXT NOT NULL,
		notification_value TEXT,
		post_id INTEGER,
		post_title TEXT,
		viewed TEXT,
		user_like_id INTEGER,
		comment_id INTEGER,
		posting_time TEXT,
		user_id INTEGER,
		CONSTRAINT fk_comment FOREIGN KEY(comment_id) REFERENCES comment(id) ON DELETE CASCADE,
		CONSTRAINT fk_user_like FOREIGN KEY(user_like_id) REFERENCES user_like(id) ON DELETE CASCADE,
		CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
		CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE	
);

CREATE TRIGGER IF NOT EXISTS log_user_like_to_activity
	AFTER INSERT ON userlike
BEGIN
DELETE from user_activity WHERE post_id = new.post_id AND notification_type = 'LIKE_DISLIKE';
	insert into user_activity (
		notification_type ,
		notification_value,
		post_id,
		post_title,
		viewed,
		user_like_id,
		comment_id,
		posting_time,
		user_id
	)
	VALUES (
		'LIKE_DISLIKE',
		new.mark,
		new.post_id,
		'',
		'false',
		new.id ,
		0,
		'',
		new.user_id
	);
END;

CREATE TRIGGER IF NOT EXISTS log_user_comment_activity 
	AFTER INSERT ON comment
BEGIN
	insert into user_activity (
		notification_type,
		notification_value,
		post_id,
		post_title,
		viewed,
		user_like_id,
		comment_id,
		posting_time,
		user_id
	)
	VALUES (
		'COMMENT',
		new.content,
		new.post_id,
		'',
		'false',
		0,
		new.id,
		new.posting_time,
		new.user_id
	);
END;

CREATE TRIGGER IF NOT EXISTS log_user_post_activity 
	AFTER INSERT ON post
BEGIN
	insert into user_activity (
		notification_type,
		notification_value,
		viewed,
		post_id,
		post_title,
		user_like_id,
		comment_id,
		posting_time,
		user_id
	)
	VALUES (
		'POST',
		new.title,
		'false',
		new.id,
		new.title,
		0,
		0,
		new.posting_time,
		new.user_id
	);
END;

-- CREATE TRIGGER IF NOT EXISTS log_edited_comment_to_activity
-- 	AFTER UPDATE ON comment
-- BEGIN
-- 	insert into user_activity (
-- 		notification_type,
-- 		notification_value,
-- 		viewed,
-- 		comment_id,
-- 		post_id,
-- 		user_id
-- 	)
-- 	VALUES (
-- 		'COMMENT_UPDATE',
-- 		new.content,
-- 		'false',
-- 		new.id,
-- 		new.post_id,
-- 		new.user_id
-- 	);
-- END;

insert into user(name,email,password, role_id) values("marat","123@mail.com","$2a$14$mVucWUpjrhpvMlNZh0aEs.KE1NOK0058Zwr7IkQuU8dnBPEYDfwt6", 3);
insert into user(name,email,password, role_id) values("nurislam","n@mail.com","$2a$14$mVucWUpjrhpvMlNZh0aEs.KE1NOK0058Zwr7IkQuU8dnBPEYDfwt6", 3);
insert into user(name,email,password, role_id) values("user","user@mail.com","$2a$14$mVucWUpjrhpvMlNZh0aEs.KE1NOK0058Zwr7IkQuU8dnBPEYDfwt6", 1);
insert into user(name,email,password, role_id) values("moderator","mod@mail.com","$2a$14$mVucWUpjrhpvMlNZh0aEs.KE1NOK0058Zwr7IkQuU8dnBPEYDfwt6", 2);

insert into userrole(name) values("user");
insert into userrole(name) values("administrator");
insert into userrole(name) values("moderator")
