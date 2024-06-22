DROP TRIGGER tr_env_ref_update_check_unique_name;
DROP TRIGGER tr_env_var_update_check_unique_name;

CREATE TRIGGER tr_env_ref_update_check_unique_name
BEFORE UPDATE ON env_ref
FOR EACH ROW
BEGIN
    SELECT RAISE(FAIL, 'name already exists in env')
    FROM vw_env_env_ref_env_var_unique_name
    WHERE env_id = NEW.env_id AND name = NEW.name AND OLD.name <> NEW.name ;
END;

CREATE TRIGGER tr_env_var_update_check_unique_name
BEFORE UPDATE ON env_var
FOR EACH ROW
BEGIN
    SELECT RAISE(FAIL, 'name already exists in env')
    FROM vw_env_env_ref_env_var_unique_name
    WHERE env_id = NEW.env_id AND name = NEW.name AND OLD.name <> NEW.name;
END;