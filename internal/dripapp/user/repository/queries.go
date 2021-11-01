package repository

const (
	getUserQuery = "select id, name, email, password, date, description from profile where email = $1;"

	getUserByIdAQuery = "select id, name, email, password, date, description from profile where id = $1;"

	createUserQuery = "INSERT into profile(email,password) VALUES($1,$2) RETURNING id, email, password;"

	updateUserQuery = `update profile set name=$1, date=$3, description=$4, imgs=$5 where email=$2 
							RETURNING id, email, password, name, email, password, date, description;`

	deleteTagsQuery = "delete from profile_tag where profile_id=$1"

	getTagsQuery = "select tag_name from tag;"

	getTagsByIdQuery = `select
							tag_name
						from
							profile p
							join profile_tag pt on(pt.profile_id = p.id)
							join tag t on(pt.tag_id = t.id)
						where
							p.id = $1;`

	getImgsByID = "SELECT imgs FROM profile WHERE id=$1;"

	insertTagsQueryFirstPart = "insert into profile_tag(profile_id, tag_id) values"
	insertTagsQueryParts     = "($1, (select id from tag where tag_name=$%d))"

	updateImgsQuery = "update profile set imgs=$2 where id=$1 returning id;"

	addReactionQuery = "insert into reactions(id1, id2, type) values ($1,$2,$3);"

	getNextUserForSwipeQuery = `select
									op.id,
									op.name,
									op.email,
									op.password,
									op.date,
									op.description
								from profile p
								join reactions r on (r.id1 = $1)
								right join profile op on (op.id = r.id2)
								left join matches m on (m.id1 = op.id)
								where
									op.name <> ''
									and op.date <> ''
									and r.id1 is null
									and op.id <> $1
									and (m.id2 is null or m.id2 <> $1)
								limit 5;`

	getUsersForMatchesQuery = `select
									op.id,
									op.name,
									op.email,
									op.password,
									op.date,
									op.description
								from profile p
								join matches m on (p.id = m.id1)
								join matches om on (om.id1 = m.id2 and om.id2 = m.id1)
								join profile op on (op.id = om.id1)
								where p.id = $1;`

	getLikesQuery = "select r.id1 from reactions r where r.id2 = $1 and r.type = 1;"

	deleteLikeQuery = "delete from reactions r where ((r.id1=$1 and r.id2=$2) or (r.id1=$2 and r.id2=$1));"

	addMatchQuery = "insert into matches(id1, id2) values ($1,$2),($2,$1);"
)
