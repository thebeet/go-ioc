CREATE TABLE `order` (
  `id` int(11) NOT NULL,
  `user_name` varchar(100) DEFAULT NULL,
  `item` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `order`
  ADD PRIMARY KEY (`id`) ,
  ADD KEY `user_name` (`user_name`);

ALTER TABLE `order`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

INSERT INTO `order` (`user_name`, `item`) VALUES
('thebeet', 'apple'), ('thebeet', 'banana'),  ('thebeet02', 'orange');