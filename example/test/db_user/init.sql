CREATE TABLE `user` (
  `id` int(11) NOT NULL,
  `name` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `user`
  ADD PRIMARY KEY (`id`) ,
  ADD KEY `name` (`name`);

ALTER TABLE `user`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

INSERT INTO `user` (`name`) VALUES
('thebeet'), ('thebeet02'),  ('thebeet03');