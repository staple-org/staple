import React, {Fragment, useEffect, useState} from "react";
import {Col, ControlLabel, FormControl, FormGroup, Grid, Row} from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import {useFormFields} from "../libs/hooksLib";
import config from "../config";
import Cookie from "js-cookie";

export default function Settings(props) {
  const [fields, handleFieldChange] = useFormFields({
    password: "",
    confirmPassword: "",
  });
  const [maxStaples, setMaxStaples] = useState(-1);
  const [isLoading, setIsLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setIsLoading(true);

    try {
      fetch(config.HOST+'/rest/api/1/user/max-staples', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get('token'),
        },
        body: JSON.stringify({
          max_staples: maxStaples,
        }),
      }).then((response) => {
        if (response.status === 200) {
          alert(`New staple count of ${maxStaples} successfully set.`);
          // window.location.reload();
          onLoad();
        } else {
          alert("Response was: " + response.statusText);
        }
      }).catch((error) => {
        alert(error.message);
      })
    } catch (e) {
      alert(e.message);
    } finally {
      setIsLoading(false);
    }
  }

  function validateChangePasswordForm() {
    return (
      fields.password.length > 0 &&
      fields.password === fields.confirmPassword
    );
  }

  async function handleChangePasswordSubmit(event) {
    event.preventDefault();
    setIsLoading(true);

    try {
      fetch(config.HOST+'/rest/api/1/user/change-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get('token'),
        },
        body: JSON.stringify({
          password: fields.password,
        }),
      }).then((response) => {
        if (response.status === 200) {
          alert("Please log in with your new password.");
          props.userHasAuthenticated(false);
          props.history.push("/login");
        } else {
          alert("Response was: " + response.statusText);
        }
      }).catch((error) => {
        alert(error.message);
      })
    } catch (e) {
      alert(e.message);
    } finally {
      setIsLoading(false);
    }
  }

  function editMaxStaples(event) {
    setMaxStaples(event.target.value)
  }

  function validateForm() {
    return true
  }

  async function onLoad() {
    try {
      fetch(config.HOST+"/rest/api/1/user/max-staples", {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get('token'),
        },
      }).then(response => response.json())
        .then(data => setMaxStaples(data.max_staples))
        .catch(e => alert(e.message));
    } catch (e) {
      alert(e);
    }
    setIsLoading(false);
  }

  useEffect(() => {
    onLoad();
  }, []);

  return (
    <Grid>
      <Row className="show-grid">
        <Col xs={12} md={8}>
          {maxStaples === -1 ? <Fragment/> : <form className="SettingsForm" onSubmit={handleSubmit}>
              <FormGroup bsSize="large" controlId="storage">
                <ControlLabel>Maximum Staples</ControlLabel>
                <FormControl
                  min="0"
                  type="number"
                  value={maxStaples}
                  onChange={editMaxStaples}
                  placeholder="Number of maximum staples to store"
                />
              </FormGroup>
              <hr/>
              <LoaderButton
                block
                type="submit"
                bsSize="large"
                isLoading={isLoading}
                disabled={!validateForm()}
              >
                Save
              </LoaderButton>
            </form>
          }
        </Col>
        <Col xs={6} md={4}>
          <form className="ChangePasswordForm" onSubmit={handleChangePasswordSubmit}>
            <FormGroup controlId="password" bsSize="large">
              <ControlLabel>New password</ControlLabel>
              <FormControl
                type="password"
                value={fields.password}
                onChange={handleFieldChange}
              />
            </FormGroup>
            <FormGroup controlId="confirmPassword" bsSize="large">
              <ControlLabel>Confirm new password</ControlLabel>
              <FormControl
                type="password"
                onChange={handleFieldChange}
                value={fields.confirmPassword}
              />
            </FormGroup>
            <hr />
            <LoaderButton
              block
              type="submit"
              bsSize="large"
              isLoading={isLoading}
              disabled={!validateChangePasswordForm()}
            >
              Change
            </LoaderButton>
          </form>
        </Col>
      </Row>
    </Grid>
  );
}