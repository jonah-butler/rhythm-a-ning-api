-- ===================================
-- associating workflows with tags  |
-- ==================================
CREATE TABLE IF NOT EXISTS workflow_tags (
    workflow_id INT NOT NULL REFERENCES workflows(workflow_id) ON DELETE CASCADE,
    tag_id      INT NOT NULL REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY     KEY (workflow_id, tag_id)
);