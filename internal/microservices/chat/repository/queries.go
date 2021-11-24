package repository

const (
	GetMessagesQuery = `select message_id, from_id, to_id, text, date
    from message
    where
      ((from_id = $1 and to_id = $2) or (from_id = $2 and to_id = $1)) and message_id < $3
    order by date
    limit 100;`

	GetLastMessageQuery = `select message_id, from_id, to_id, text, date
    from message
    where
      (from_id = $1 and to_id = $2)
      or (from_id = $2 and to_id = $1 )
    order by date desc
    limit 1;`

	SendMessageQuery = `
	insert into message(from_id, to_id, text) values ($1,$2,$3) returning message_id, from_id, to_id, text, date;
	`

	GetChatsQuery = `
	select
		op.id as FromUserID, op.name as name, op.imgs[1] as img
	from
		profile p
		join message m on p.id = m.from_id
		join profile op on op.id = m.to_id
	where m.from_id=$1
	union select
		p.id as FromUserID, p.name as name, p.imgs[1] as img
	from
		profile p
		join message m on p.id = m.from_id
		join profile op on op.id = m.to_id
	where m.to_id=$1
	group by op.id, op.name, op.imgs[1], p.id, p.id, p.name, p.imgs[1];
	`
)
