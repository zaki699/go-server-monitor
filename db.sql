create database live_encoding;

CREATE USER 'live-user'@'84.255.30.116' IDENTIFIED WITH mysql_native_password BY 'GAYxmZfFkB.Twlc7E!';
GRANT SELECT, INSERT, UPDATE ON live_encoding.* TO 'live-user'@'84.255.30.116';
FLUSH PRIVILEGES;

CREATE TABLE `sessions` (
  `session_id` int NOT NULL AUTO_INCREMENT,
  `channel_name` char(50) NOT NULL,
  `hostname` char(50) NOT NULL,
  `preset` char(50) NOT NULL,
  `name` char(150) NOT NULL,
  `definition` char(3) NOT NULL,
  `codec` varchar(5) NOT NULL,
  `optimizer_enabled` tinyint(1) NOT NULL DEFAULT '0',
  `status` enum('running','end') NOT NULL DEFAULT 'running',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ended_at` timestamp NULL DEFAULT NULL,
  `cmd` varchar(5000) DEFAULT NULL,
  PRIMARY KEY (`session_id`)
) ENGINE=InnoDB

CREATE TABLE stats (
    stat_id INT AUTO_INCREMENT PRIMARY KEY,
    frames INT NOT NULL DEFAULT 0,
    drop_frames INT NOT NULL DEFAULT 0,
    dup_frames INT NOT NULL DEFAULT 0,
    session_id INT NOT NULL,
    speed FLOAT NOT NULL DEFAULT 0.0,
    bitrate CHAR(10) NOT NULL DEFAULT "N/A",
    encoding_time BIGINT NOT NULL DEFAULT 0,
    streams_qp CHAR(50) NOT NULL ,
    cmd VARCHAR(2000) NOT NULL DEFAULT "N/A",
    fps FLOAT NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id)
);

CREATE TABLE stats_aggregated (
    stat_aggregated_id INT AUTO_INCREMENT PRIMARY KEY,
    nbr_frames_encoded INT NOT NULL DEFAULT 0,
    avg_drop_frames INT NOT NULL DEFAULT 0,
    avg_dup_frames INT NOT NULL DEFAULT 0,
    session_id INT NOT NULL,
    avg_speed FLOAT NOT NULL DEFAULT 0.0,
    avg_bitrate CHAR(10) NOT NULL DEFAULT "N/A",
    encoding_time BIGINT NOT NULL DEFAULT 0,
    avg_streams_qp CHAR(50) NOT NULL ,
    avg_fps FLOAT NOT NULL DEFAULT 0.0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id)
);

CREATE TABLE logs (
    log_id INT AUTO_INCREMENT PRIMARY KEY,
    level CHAR(20) NOT NULL,
    module CHAR(20) NOT NULL,
    session_id INT NOT NULL,
    message VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id)
);

