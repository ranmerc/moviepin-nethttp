INSERT INTO users (user_id, username, email, password)
VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'dummyuser1', 'dummy@email.com', 'dummyhash1');

INSERT INTO movies (movie_id, title, release_date, genre, director, description)
VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'Dummy Movie 1', '2021-01-01', 'Action', 'Dummy Director 1', 'A dummy description 1.'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Dummy Movie 2', '2021-01-02', 'Comedy', 'Dummy Director 2', 'A dummy description 2.');


INSERT INTO reviews (user_id, movie_id, rating, review_text)
VALUES 
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 5, 'A dummy review text 1.'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 5, 'A dummy review text 1.'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 4, 'A dummy review text 2.');