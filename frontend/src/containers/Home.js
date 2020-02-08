import React, {useEffect, useState} from "react";
import {Col, FormGroup, Grid, ListGroup, ListGroupItem, PageHeader, Row} from "react-bootstrap";
import {LinkContainer} from "react-router-bootstrap";
import "./Home.css";
import config from "../config";
import Cookie from "js-cookie";
import LoaderButton from "../components/LoaderButton";

export default function Home(props) {
  const [staples, setStaples] = useState([]);
  const [nextStaple, setNextStaple] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const [isArchiving, setIsArchiving] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  function handleArchive(event, id) {
    setIsLoading(true);
    event.preventDefault();
    try {
      fetch(config.HOST+`/rest/api/1/staple/${id}/archive`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get('token'),
        },
      }).then(response => {
        if (response.ok) {
          setIsArchiving(false);
          props.history.push('/')
        }
      }).catch(e => alert(e.message));
    } catch (e) {
      alert(e);
    }
  }

  function handleDelete(event, id) {
    setIsDeleting(true);
    event.preventDefault();
    try {
      fetch(config.HOST+`/rest/api/1/staple/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get('token'),
        },
      }).then(response => {
        if (response.ok) {
          setIsDeleting(false);
          props.history.push('/')
        }
      }).catch(e => alert(e.message));
    } catch (e) {
      alert(e);
    }
  }

  function renderStapleView(astaples) {
    let list = undefined;
    if (astaples.staples) {
      list = astaples.staples.map(staple => {
        return (
          <div key={staple.id} className="staple-list">
            <strong>{"Name: " + staple.name}</strong><br/>
            {"Created: " + new Date(staple.created_at).toLocaleString()}
          </div>
        )
      })
    }
    let s = undefined;
    if (nextStaple.staple) {
      s = (
        <div className="Staples">
          <form>
            <FormGroup controlId="content" className="staple-view">
              <pre>{nextStaple.staple.content}</pre>
            </FormGroup>
            <LoaderButton
              block
              bsSize="large"
              bsStyle="info"
              onClick={(event) => handleArchive(event, nextStaple.staple.id)}
              isLoading={isArchiving}
              className="archive-button"
            >
              Archive
            </LoaderButton>
            <LoaderButton
              block
              bsSize="large"
              bsStyle="info"
              onClick={(event) => handleDelete(event, nextStaple.staple.id)}
              isLoading={isDeleting}
              className="delete-button"
            >
              Delete
            </LoaderButton>
          </form>
          <br/>
        </div>
      )
    }
    return (
      <div>
        <Grid>
          <Row className="show-grid">
            <Col xs={12} md={8}>
              {s}
            </Col>
            <Col xs={6} md={4}>
              {list}
            </Col>
          </Row>
        </Grid>
        <LinkContainer key="new" to="/staples/new" className="create-container">
          <ListGroupItem>
            <h4>
              <b>{"\uFF0B"}</b> Create a new staple
            </h4>
          </ListGroupItem>
        </LinkContainer>
      </div>
    )
  }

  function renderLander() {
    return (
      <div className="lander">
        <h1>Staple</h1>
        <p>A stack based bookmark.</p>
      </div>
    );
  }

  function renderStaples() {
    return (
      <div className="staples">
        <PageHeader>Your Staples</PageHeader>
        <ListGroup>
          {!isLoading && renderStapleView(staples)}
        </ListGroup>
      </div>
    );
  }

  useEffect(() => {
    async function onLoad() {
      if (!props.isAuthenticated) {
        return;
      }

      try {
        fetch(config.HOST+"/rest/api/1/staple", {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + Cookie.get('token'),
          },
        }).then(response => response.json())
          .then(data => setStaples(data))
          .catch(e => alert(e.message));
      } catch (e) {
        alert(e);
      }

      try {
        fetch(config.HOST+"/rest/api/1/staple/next", {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + Cookie.get('token'),
          },
        }).then(response => response.json())
          .then(data => setNextStaple(data))
          .catch(e => alert(e.message));
      } catch (e) {
        alert(e);
      }

      setIsLoading(false);
    }

    onLoad();
  }, [props.isAuthenticated, isLoading, isDeleting]);

  return (
    <div className="Home">
      {props.isAuthenticated ? renderStaples() : renderLander()}
    </div>
  );
}