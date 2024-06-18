use actix_web::http::StatusCode;
use std::process::Command;
use std::{fs, str};

struct UserShadowData {
    hash: String,
    hash_type: String,
    salt: String,
}

pub struct UserPasswdData {
    pub username: String,
    pub password: String,
    pub uid: u32,
    pub gid: u32,
    pub gecos: String,
    pub home: String,
    pub shell: String,
}

pub struct UserServiceResult {
    pub message: String,
    pub status: StatusCode,
}

pub struct UserServiceError {
    pub message: String,
    pub status: StatusCode,
}

fn ok() -> UserServiceResult {
    UserServiceResult {
        message: "Ok".to_owned(),
        status: StatusCode::OK,
    }
}

fn created() -> UserServiceResult {
    UserServiceResult {
        message: "Created".to_owned(),
        status: StatusCode::CREATED,
    }
}

fn unauthorized() -> UserServiceError {
    UserServiceError {
        message: "Unauthorized".to_owned(),
        status: StatusCode::UNAUTHORIZED,
    }
}

fn forbidden() -> UserServiceError {
    UserServiceError {
        message: "Forbidden".to_owned(),
        status: StatusCode::FORBIDDEN,
    }
}

fn not_found() -> UserServiceError {
    UserServiceError {
        message: "Not found".to_owned(),
        status: StatusCode::NOT_FOUND,
    }
}

fn internal_error() -> UserServiceError {
    UserServiceError {
        message: "Internal server error".to_owned(),
        status: StatusCode::INTERNAL_SERVER_ERROR,
    }
}

fn find_shadow_entry(username: &str) -> Result<String, UserServiceError> {
    for line in fs::read_to_string("/etc/shadow").unwrap().lines() {
        if line.starts_with(username) {
            return Ok(line.to_owned());
        }
    }

    return Err(not_found());
}

fn find_passwd_entry(username: &str) -> Result<String, UserServiceError> {
    for line in fs::read_to_string("/etc/passwd").unwrap().lines() {
        if line.starts_with(username) {
            return Ok(line.to_owned());
        }
    }

    return Err(not_found());
}

fn parse_shadow(shadow_entry: &str) -> Result<UserShadowData, UserServiceError> {
    let split: Vec<&str> = shadow_entry.split(":").collect();
    match split[..] {
        [_, password_hash, ..] => {
            let split: Vec<&str> = password_hash.split("$").collect();
            match split[..] {
                [_, "y", ..] => {
                    log::info!("Illegal algorithm: y");
                    Err(forbidden())
                }
                [_, hash_type, salt, ..] => Ok(UserShadowData {
                    hash: password_hash.to_owned(),
                    hash_type: hash_type.to_owned(),
                    salt: salt.to_owned(),
                }),
                [..] => {
                    log::info!("Invalid hash entry");
                    Err(forbidden())
                }
            }
        }
        [..] => Err(forbidden()),
    }
}

fn parse_passwd(passwd_entry: &str) -> Result<UserPasswdData, UserServiceError> {
    let split: Vec<&str> = passwd_entry.split(":").collect();
    match split[..] {
        [username, password, uid, gid, gecos, home, shell, ..] => {
            let uid = uid.parse::<u32>().map_err(|_| forbidden())?;
            let gid = gid.parse::<u32>().map_err(|_| forbidden())?;
            Ok(UserPasswdData {
                username: username.to_owned(),
                password: password.to_owned(),
                uid,
                gid,
                gecos: gecos.to_owned(),
                home: home.to_owned(),
                shell: shell.to_owned(),
            })
        }
        [..] => Err(forbidden()),
    }
}

fn validate_user_password(shadow: UserShadowData, password: &str) -> Result<(), UserServiceError> {
    let result = Command::new("openssl")
        .arg("passwd")
        .arg(format!("-{}", shadow.hash_type))
        .arg("-salt")
        .arg(shadow.salt)
        .arg(password)
        .output()
        .map_err(|_| internal_error())?;

    if !result.status.success() {
        return Err(internal_error());
    }

    let hash = str::from_utf8(result.stdout.as_slice()).map_err(|_| internal_error())?;

    if hash.trim() == shadow.hash {
        Ok(())
    } else {
        Err(unauthorized())
    }
}

fn create_user(username: &str, password: &str) -> Result<(), UserServiceError> {
    let status = Command::new("adduser")
        .arg("-D")
        .arg(username)
        .arg("-s")
        .arg("/bin/zsh")
        .status()
        .map_err(|_| internal_error())?;

    if !status.success() {
        return Err(internal_error());
    }

    let status = Command::new("sh")
        .arg("-c")
        .arg(format!("echo {}:{} | chpasswd", username, password))
        .status()
        .map_err(|_| internal_error())?;

    if status.success() {
        Ok(())
    } else {
        Err(internal_error())
    }
}

pub fn get_user_data(username: &str) -> Result<UserPasswdData, UserServiceError> {
    let passwd_entry = find_passwd_entry(username)?;
    parse_passwd(&passwd_entry)
}

pub fn create(username: &str, password: &str) -> Result<UserServiceResult, UserServiceError> {
    match find_shadow_entry(username) {
        Ok(shadow_entry) => {
            let shadow = parse_shadow(&shadow_entry)?;
            validate_user_password(shadow, password).and_then(|_| Ok(ok()))
        }
        Err(_) => create_user(username, password).and_then(|_| Ok(created())),
    }
}

pub fn login(username: &str, password: &str) -> Result<UserServiceResult, UserServiceError> {
    let shadow_entry = find_shadow_entry(username)?;
    let shadow = parse_shadow(&shadow_entry)?;
    validate_user_password(shadow, password).and_then(|_| Ok(ok()))
}
