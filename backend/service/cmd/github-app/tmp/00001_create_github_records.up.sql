CREATE TABLE bounty_creator (
    id SERIAL,
    entity_id INT NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bounty_status (
    id SERIAL,
    name VARCHAR(255) NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/*
    github_bounty records the bounties created by a bounty creator.
*/
CREATE TABLE bounty (
    id SERIAL,
    entity_id INT NOT NULL,
    url VARCHAR(255) NOT NULL,
    issue_id INT NOT NULL PRIMARY KEY,
    repo_id INT NOT NULL,
    repo_name VARCHAR(255) NOT NULL,
    owner_id INT NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT bounty_creator
        FOREIGN KEY(entity_id)
            REFERENCES bounty_creator(entity_id),
    CONSTRAINT bounty_status
        FOREIGN KEY(status)
            REFERENCES bounty_status(name)
);

INSERT INTO bounty_status (name) VALUES ('open');
INSERT INTO bounty_status (name) VALUES ('closed');
INSERT INTO bounty_status (name) VALUES ('complete');
