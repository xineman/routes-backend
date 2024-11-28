create table tracks (
  id          serial primary key,
  created_at  timestamp not null default now(),
  updated_at  timestamp not null default now(),
  name        text not null,
  file_name   text not null UNIQUE
)
