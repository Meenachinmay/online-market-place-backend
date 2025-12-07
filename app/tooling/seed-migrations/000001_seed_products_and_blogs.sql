-- Products
INSERT INTO products (id, name, description, price, buyer_reward_points, author_reward_points) VALUES
('68d4d95a-5df3-40c9-b0e2-fb5c30904e09', 'Soda Air Max 1', 'Classic comfort and style for everyday wear.', 100, 10, 5),
('07ea4a41-ea43-465b-a121-cedc5becf52e', 'Soda Runner Low', 'Lightweight running shoes for speed.', 200, 20, 10),
('a3a35d83-fe83-45a1-b6f1-fb6f03debf84', 'Soda High Tops', 'Vintage style high tops with modern grip.', 300, 30, 15),
('719b2afe-034d-4ec2-8b51-aae6fced2884', 'Soda Court Vision', 'Basketball inspired sneakers for the court.', 400, 40, 20),
('126f3a3e-302d-4286-8971-c3e98f5778ef', 'Soda Street King', 'Urban street style with premium leather.', 500, 50, 25),
('d142a19c-6e05-4183-9c48-99330d82d54d', 'Soda Canvas Slip-On', 'Casual slip-ons for a relaxed vibe.', 600, 60, 30),
('f1fa509f-57a5-4d28-a111-863d470978da', 'Soda Trail Blazer', 'Rugged outsole for off-road adventures.', 700, 70, 35),
('25d45deb-c5ab-449d-b9a8-5f01d061972e', 'Soda Retro 90s', 'Throwback design with chunky soles.', 800, 80, 40),
('54972b8c-b1fd-4cab-8788-53a5fd193a1f', 'Soda Knit Runner', 'Breathable knit upper for maximum airflow.', 900, 90, 45),
('61aad777-8e7c-4135-ab2b-590561bcaea5', 'Soda Pro Skater', 'Durable suede reinforced for skating.', 1000, 100, 50),
('2db45c6f-e98a-4d75-a3de-bd39a4bd4542', 'Soda Elite Racer', 'Professional grade marathon shoes.', 1100, 110, 55),
('4852788b-47eb-4ebf-8699-13b8ad58ff76', 'Soda Limited Edition', 'Exclusive colorway, limited stock.', 1200, 120, 60),
('04777eaa-7e1a-4418-88ae-1eda740b3c3d', 'Soda Tech Future', 'Futuristic design with auto-lacing tech.', 1300, 130, 65),
('838d7318-5ea5-444e-bfd1-bec8e362617e', 'Soda Classic White', 'Minimalist white sneakers that go with everything.', 1400, 140, 70),
('257270b6-f9e1-40af-80d2-e11980800e58', 'Soda Midnight Black', 'Sleek all-black design for night outs.', 1500, 150, 75)
ON CONFLICT (id) DO NOTHING;

-- Blogs
INSERT INTO blogs (id, author_id, content, product_id) VALUES
('c0fa0d06-156a-4ddf-bfe1-e048b5597dba', '91183969-4b3a-4692-8ec6-74e2caf131af', 'These Soda Air Max 1s are incredibly comfortable for walking all day. Highly recommend!', '68d4d95a-5df3-40c9-b0e2-fb5c30904e09'),
('cbd832e2-2957-4492-b7d7-b63e8c31711a', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Just got my Soda Runner Lows. They are so light, it feels like running on clouds.', '07ea4a41-ea43-465b-a121-cedc5becf52e'),
('610c6ca6-ec21-47e9-8831-b78cc209c124', '91183969-4b3a-4692-8ec6-74e2caf131af', 'The vintage look on these Soda High Tops is fire. Great grip too.', 'a3a35d83-fe83-45a1-b6f1-fb6f03debf84'),
('d4bc64a4-26de-464f-8cb1-df29ac3c0355', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Soda Court Vision has improved my game. Excellent ankle support.', '719b2afe-034d-4ec2-8b51-aae6fced2884'),
('be6714ee-49f9-45cb-bce4-0e8bea8b14ce', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Rocking the Soda Street Kings today. The leather quality is premium.', '126f3a3e-302d-4286-8971-c3e98f5778ef'),
('82ce1b67-f63a-47dd-b105-79243a4e9774', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Soda Canvas Slip-Ons are my go-to for quick errands. Super easy.', 'd142a19c-6e05-4183-9c48-99330d82d54d'),
('76f3f15a-de44-4231-9b53-a6e1d4ec542b', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Took the Soda Trail Blazers hiking this weekend. No slips, great traction.', 'f1fa509f-57a5-4d28-a111-863d470978da'),
('030fca7b-f3fe-4638-a7bf-364849daa476', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Loving the chunky sole on these Soda Retro 90s. Total nostalgia trip.', '25d45deb-c5ab-449d-b9a8-5f01d061972e'),
('8b259397-e77e-473f-a705-544cccc84aa1', '91183969-4b3a-4692-8ec6-74e2caf131af', 'My feet breathe so well in the Soda Knit Runners. Perfect for summer.', '54972b8c-b1fd-4cab-8788-53a5fd193a1f'),
('f8d66b18-218a-4edc-80f3-0910e0c02996', '91183969-4b3a-4692-8ec6-74e2caf131af', 'Soda Pro Skaters holding up well after a week of intense sessions.', '61aad777-8e7c-4135-ab2b-590561bcaea5')
ON CONFLICT (id) DO NOTHING;

