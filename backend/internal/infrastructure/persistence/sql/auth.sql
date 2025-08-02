-- name: GetUserFamilyIds :many
SELECT fm.family_id FROM family_members fm
where fm.user_id = $1;