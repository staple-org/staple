import React, {useState} from "react";
import "./ResetPassword.css";
import config from "../config";
import {ControlLabel, FormControl, FormGroup, HelpBlock} from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";

export default function ResetPassword(props) {
  const [isVerifying, setIsVerifying] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');

  function validateForm() {
    return email.length > 0;
  }

  function validateConfirmationForm() {
    return code.length > 0;
  }

  async function handleSubmit(event) {
    event.preventDefault();

    setIsLoading(true);

    try {
      fetch(config.HOST+'/reset', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
        }),
      }).then((response) => {
        if (response.status === 200) {
          setIsVerifying(true);
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

  async function handleConfirmationSubmit(event) {
    event.preventDefault();

    setIsLoading(true);

    try {
      fetch(config.HOST+'/verify', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
          code: code,
        }),
      }).then((response) => {
        if (response.status === 200) {
          alert("Please log in with your new password.");
          setIsVerifying(false);
          props.history.push("/");
        } else {
          alert("Confirmation code did not match.");
        }
      }).catch((error) => {
        alert(error.message);
      });
    } catch (e) {
      alert(e.message);
    } finally {
      setIsLoading(false);
    }
  }

  function renderConfirmationForm() {
    return (
      <form onSubmit={handleConfirmationSubmit}>
        <FormGroup controlId="confirmationCode" bsSize="large">
          <ControlLabel>Confirmation Code</ControlLabel>
          <FormControl
            autoFocus
            type="text"
            onChange={setCodeField}
            value={code}
          />
          <HelpBlock>Please check your email for the code.</HelpBlock>
        </FormGroup>
        <LoaderButton
          block
          type="submit"
          bsSize="large"
          isLoading={isLoading}
          disabled={!validateConfirmationForm()}
        >
          Verify
        </LoaderButton>
      </form>
    );
  }

  function setEmailField(event) {
    setEmail(event.target.value)
  }

  function setCodeField(event) {
    setCode(event.target.value)
  }

  function renderForm() {
    return (
      <form onSubmit={handleSubmit}>
        <FormGroup controlId="email" bsSize="large">
          <ControlLabel>Email</ControlLabel>
          <FormControl
            autoFocus
            type="email"
            value={email}
            onChange={setEmailField}
          />
        </FormGroup>
        <LoaderButton
          block
          type="submit"
          bsSize="large"
          isLoading={isLoading}
          disabled={!validateForm()}
        >
          Reset Password
        </LoaderButton>
      </form>
    );
  }
  return (
    <div className="ResetPassword">
      {isVerifying ? renderConfirmationForm() : renderForm()}
    </div>
  );
}