import React, { useState } from "react";
import { FormGroup, FormControl } from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import "./AddStaple.css";
import Cookie from "js-cookie";

export default function AddStaple(props) {
  const [content, setContent] = useState("");
  const [title, setTitle] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  function validateForm() {
    return content.length > 0 && title.length > 0;
  }


  async function handleSubmit(event) {
    // Call the api to create a staple.
    event.preventDefault();
    setIsLoading(true);
    try {
      fetch("/rest/api/1/staple", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get("token"),
        },
        body: JSON.stringify({
          name: title,
          content: content,
        }),
      }).then(response => {
        if (response.ok) {
          setIsLoading(false);
          props.history.push("/");
        } else {
          setIsLoading(false);
          alert("Not allowed to create more staples");
        }
      }).catch(e => alert(e.message))
    } catch (e) {
      alert(e);
      setIsLoading(false);
    }
  }

  return (
    <div className="AddStaple">
      <form onSubmit={handleSubmit}>
        <FormGroup controlId="Title">
          <FormControl type="text" value={title} onChange={e => setTitle(e.target.value)}/>
        </FormGroup>
        <FormGroup controlId="content">
          <FormControl
            value={content}
            componentClass="textarea"
            onChange={e => setContent(e.target.value)}
            rows={30}
            scrolling
          />
        </FormGroup>
        <LoaderButton
          block
          type="submit"
          bsSize="large"
          bsStyle="primary"
          isLoading={isLoading}
          disabled={!validateForm()}
        >
          Create
        </LoaderButton>
      </form>
    </div>
  );
}