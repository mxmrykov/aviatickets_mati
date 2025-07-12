SELECT id, telegram_id, username, first_name, last_name, created_at::text
FROM users WHERE telegram_id = @telegram_id