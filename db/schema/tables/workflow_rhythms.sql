-- ===================================================
-- table for rhythms stored within a workflow        |
-- references its workflow, and rhythm               |
-- with the addition of position within the workflow |
-- ===================================================
CREATE TABLE IF NOT EXISTS workflow_rhythms (
    workflow_rhythm_id SERIAL PRIMARY KEY,
    workflow_id        INT NOT NULL REFERENCES workflows(workflow_id) ON DELETE CASCADE,
    rhythm_id          INT NOT NULL REFERENCES rhythms(rhythm_id) ON DELETE CASCADE,
    measures           SMALLINT NOT NULL,
    position           SMALLINT NOT NULL
);