-- Rhythm tag insert
INSERT INTO tags (name)
VALUES
  ('Spain')
ON CONFLICT (name)
DO NOTHING;

-- Spain related rhythms

-- Buleria 1
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    210,
    6,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'duplet'),
    ARRAY[1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0]::SMALLINT[],
    false,
    6,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    'Bulería [1]',
    'A variant of the 12 beat rhythm popular in flamenco music that can be thought of as a hemiola where a measure of 6/8 is followed by a measure of 3/4. Often the first beat is felt as the 12th note of the phrase.',
    (SELECT tag_id FROM tags WHERE name = 'Spain')
  )

-- Buleria 2
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    210,
    6,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'duplet'),
    ARRAY[1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 1, 0]::SMALLINT[],
    false,
    6,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    'Bulería [2]',
    'A variant of the 12 beat rhythm popular in flamenco music where the third accent is shifted forward by a single 8th note. Often the first beat is felt as the 12th note of the phrase.',
    (SELECT tag_id FROM tags WHERE name = 'Spain')
  )

-- Siguiriyas
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    150,
    6,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'duplet'),
    ARRAY[1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0]::SMALLINT[],
    false,
    6,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    'Siguiriyas',
    'A common 12 note common used in Spanish flamencos. This pattern contrasts rhythms like Buleria but can be felt as a Bulerias starting at the midpoint.',
    (SELECT tag_id FROM tags WHERE name = 'Spain')
  )

-- Caribbean related rhythms

-- Rhythm tag insert
INSERT INTO tags (name)
VALUES
  ('Caribbean')
ON CONFLICT (name)
DO NOTHING;

-- 2:3 Son Clave
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    120,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'quadruplet'),
    ARRAY[0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0]::SMALLINT[],
    false,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    '2:3 Son Clave',
    'A foundational two-measure rhythm found in Afro-Cuban and Latin music where the first beat of the 2 beat partial does not begin on the 1. Part of the son clave rhythms which are considered to be more traditional and straighter than the rumba family of claves.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

-- 3:2 Son Clave
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    120,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'quadruplet'),
    ARRAY[1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0]::SMALLINT[],
    false,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    '2:3 Son Clave',
    'A foundational two-measure rhythm found in Afro-Cuban and Latin music where the starting point is felt on the 1. Part of the son clave rhythms which are considered to be more traditional and straighter than the rumba family of claves.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

-- 2:3 Rumba Clave
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    120,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'quadruplet'),
    ARRAY[0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1]::SMALLINT[],
    false,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    '2:3 Rumba Clave',
    'A foundational two-measure rhythm found in Afro-Cuban and Latin music. Part of the rumba clave rhythms which are considered more syncopated, giving it a more distinct and rolling sound.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

-- 3:2 Rumba Clave
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    120,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'quadruplet'),
    ARRAY[1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0]::SMALLINT[],
    false,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'intermediate'),
    '3:2 Rumba Clave',
    'A foundational two-measure rhythm found in Afro-Cuban and Latin music. Part of the rumba clave rhythms which are considered more syncopated, giving it a more distinct and rolling sound.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

-- Tresillo
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    160,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'duplet'),
    ARRAY[1, 0, 0, 1, 0, 0, 1, 0]::SMALLINT[],
    false,
    4,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1, 1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'beginner'),
    'Tresillo',
    'Perhaps the most fundamental rhythmic pattern in Afro-Cuban and Latin-American music. It''s influence has reached across the world, helping inspire other musical forms found in New Orleans, Jamaica and Brazil.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

-- Habanera
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    110,
    2,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'duplet'),
    ARRAY[1, 0, 0, 1, 1, 0, 1, 0]::SMALLINT[],
    false,
    2,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'beginner'),
    'Habanera',
    'A rhythm sibling and more simplified version of the tresillo. Otherwise known as the contradanza, it is a Spanish and Spanish-American version of the contradanse. In Cuba it became the first written music to be rhythmically based on an African rhythm.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

  -- Rhythm tag insert
INSERT INTO tags (name)
VALUES
  ('Polyrhythm')
ON CONFLICT (name)
DO NOTHING;

-- Habanera
INSERT INTO rhythms
  (
    bpm,
    beats,
    subdivision,
    state,
    is_poly,
    poly_beats,
    poly_subdivision,
    poly_state,
    level_id,
    name,
    description,
    tag_id,
  )
VALUES
  (
    110,
    2,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'duplet'),
    ARRAY[1, 0, 0, 1, 1, 0, 1, 0]::SMALLINT[],
    false,
    2,
    (SELECT subdivision_id FROM subdivision_types WHERE name = 'base'),
    ARRAY[1, 1]::SMALLINT[],
    (SELECT level_id FROM levels WHERE name = 'beginner'),
    'Habanera',
    'A rhythm sibling and more simplified version of the tresillo. Otherwise known as the contradanza, it is a Spanish and Spanish-American version of the contradanse. In Cuba it became the first written music to be rhythmically based on an African rhythm.',
    (SELECT tag_id FROM tags WHERE name = 'Caribbean')
  )

