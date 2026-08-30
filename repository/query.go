package repository

import (
	_ "embed"
)

// RHYTHM QUERIES
//
//go:embed subdivisions/get_subdivision_types.sql
var GET_SUBDIVISION_TYPES string

//go:embed levels/get_rhythm_levels.sql
var GET_RHYTHM_LEVELS string

//go:embed rhythms/create_rhythm.sql
var CREATE_RHYTHM string

//go:embed rhythms/get_rhythm_by_id.sql
var GET_RHYTHM_BY_ID string

//go:embed rhythms/update_rhythm_by_id.sql
var UPDATE_RHYTHM_BY_ID string

//go:embed rhythms/delete_rhythm_by_id.sql
var DELETE_RHYTHM_BY_ID string

//go:embed rhythms/get_rhythms.sql
var GET_RHYTHMS string

//go:embed rhythms/insert_rhythm_tags.sql
var INSERT_RHYTHM_TAGS string

// USER QUERIES
//
//go:embed users/register_user.sql
var REGISTER_USER string

//go:embed users/insert_user_registration_hash.sql
var INSERT_USER_REGISTRATION_HASH string

//go:embed users/update_user_registration_hash.sql
var UPDATE_USER_REGISTRATION_HASH string

//go:embed users/get_user_id_by_username_and_email.sql
var GET_USERID_BY_USERNAME_AND_EMAIL string

//go:embed users/get_user_by_hash.sql
var GET_USER_BY_HASH string

//go:embed users/delete_user_registration_hash.sql
var DELETE_USER_REGISTRATION_HASH string

//go:embed users/authenticate_user.sql
var AUTHENTICATE_USER string

//go:embed users/insert_refresh_token.sql
var INSERT_REFRESH_TOKEN string

//go:embed users/get_user_by_id.sql
var GET_USER_BY_ID string

//go:embed users/insert_user_token.sql
var INSERT_USER_TOKEN string

//go:embed users/update_user_password.sql
var UPDATE_USER_PASSWORD string

//go:embed users/delete_user_token.sql
var DELETE_USER_TOKEN string

//go:embed users/get_user_token.sql
var GET_USER_TOKEN string
