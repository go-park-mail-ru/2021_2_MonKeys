package repository

const (
	GetMessages = `select message_id, from_id, to_id, text, date
    from message
    where
      ((from_id = $1 and to_id = $2) or (from_id = $2 and to_id = $1)) and message_id < $3
    order by date
    limit 100;`

	GetLastMessage = `select message_id, from_id, to_id, text, date
    from message
    where
      (from_id = $1 and to_id = $2)
      or (from_id = $2 and to_id = $1 )
    order by date desc
    limit 1;`

	SendNessage = `
	insert into message(from_id, to_id, text) values ($1,$2,$3) returning message_id, from_id, to_id, text, date;
	`

	InitChat = `
	insert into message(from_id, to_id) values ($1,$2);
	insert into message(from_id, to_id) values ($2,$1);
	`

	GetChats = `
	select
		p.id as FromUserID, p.name, p.imgs[1] as img
	from
		profile p
		join message m on p.id = m.from_id
		join profile op on op.id = m.to_id
	where (m.from_id=$1 or m.to_id=$1) and (p.id<>$1)
	group by p.id, p.name, p.imgs[1];
	`

	// GetChats = `
	// select
	// 	(case when p.id=$1 then op.id else p.id end) as FromUserID, p.name, p.imgs[1]as img
	// from
	// 	profile p
	// 	join message m on p.id = m.from_id
	// 	join profile op on op.id = m.to_id
	// where (m.from_id=$1 or m.to_id=$1)
	// group by p.id, p.name, p.imgs[1], op.id
	// having count(*)=1;
	// `
)
