USE wallet_service;

CREATE TABLE users (
  id          INT NOT NULL AUTO_INCREMENT,
  username    VARCHAR(255) NOT NULL UNIQUE,
  email       VARCHAR(255) NOT NULL UNIQUE,
  name        VARCHAR(255) NOT NULL DEFAULT "",
  password    VARCHAR(255) NOT NULL,
  ip_address  VARCHAR(50) NOT NULL DEFAULT "",
  active      TINYINT(1) NOT NULL DEFAULT 1,

  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE roles (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL UNIQUE,

  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE user_role (
  user_id INT NOT NULL,
  role_id INT NOT NULL,

  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (role_id) REFERENCES roles(id),
  CONSTRAINT unq_user_role UNIQUE (user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE permissions (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL UNIQUE,

  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE role_permission (
  role_id INT NOT NULL,
  permission_id INT NOT NULL,

  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (permission_id) REFERENCES permissions(id),
  CONSTRAINT unq_role_permission UNIQUE (role_id, permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `permissions` (`id`, `name`) VALUES
-- backend api names --
(1, 'listusers'),
(2, 'createuser'),
(3, 'updateuser'),
(4, 'deactivateuser'),
(5, 'activateuser'),
(6, 'createuserrole'),
(7, 'deleteuserrole'),
(8, 'listroles'),
(9, 'createrole'),
(10, 'updaterole'),
(11, 'deleterole'),
(12, 'createrolepermission'),
(13, 'deleterolepermission'),
(14, 'listpermissions'),
(15, 'createpermission'),
(16, 'updatepermission'),
(17, 'deletepermission'),
(18, 'getblockcount'),
(19, 'getblockcountbysymbol'),
(20, 'gethealthcheck'),
(21, 'gethealthcheckbysymbol'),
(22, 'getlog'),
(23, 'updatemaintlist'),

-- frontend page names --
(23, 'UserList'),
(24, 'RoleList'),
(25, 'PermissionList'),
(26, 'NodesInfo');

INSERT INTO `roles` (`id`, `name`) VALUES
(1, 'admin'),
(2, 'user');

-- Add admin role to all permissions
INSERT INTO `role_permission` (`role_id`, `permission_id`) VALUES
(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(1, 6),
(1, 7),
(1, 8),
(1, 9),
(1, 10),
(1, 11),
(1, 12),
(1, 13),
(1, 14),
(1, 15),
(1, 16),
(1, 17),
(1, 18),
(1, 19),
(1, 20),
(1, 21),
(1, 22),
(1, 23),
(1, 24),
(1, 25),
(1, 26);

-- password is hashed golang/bcrypt value of '12345678'
INSERT INTO `users` (`id`, `username`, `name`, `email`, `password`, `ip_address`) VALUES
(1, 'admin', 'Admin', 'admin@admin.com', '$2a$10$xNGpvGcIqbFtFpzSLmxCReQDRYWlKWitdHn0naC5wiz24dumti.gW', '');

INSERT INTO `user_role` (`user_id`, `role_id`) VALUES
(1, 1);
