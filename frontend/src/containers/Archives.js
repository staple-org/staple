import React, { useState, useEffect } from "react";
import config from "../config";
import Cookie from "js-cookie";
import {FormGroup, Grid, Row, Col, PageHeader, ListGroup} from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import "./Archives.css";

export default function Archives(props) {
  const [staple, setStaple] = useState({});
  const [staples, setStaples] = useState([]);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    async function onLoad() {
      setIsLoading(true);
      try {
        fetch(config.HOST+"/rest/api/1/staple/archive", {
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
      setIsLoading(false);
    }
    onLoad();
  }, [props.match.params.id]);

  async function handleDelete(id) {
    setIsDeleting(true);
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
          setStaple({});
          window.location.reload();
        }
      }).catch(e => alert(e.message));
    } catch (e) {
      alert(e);
    } finally {
      setIsDeleting(false);
    }
  }

  function fetchStaple(id) {
    try {
      fetch(config.HOST+`/rest/api/1/staple/${id}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + Cookie.get('token'),
        },
      }).then(response => response.json())
        .then(data => setStaple(data))
        .catch(e => alert(e.message));
    } catch (e) {
      alert(e);
    }
  }

  function renderArchiveView(astaples) {
    let list = undefined;
    if (astaples.staples) {
      list = astaples.staples.map(staple => {
        return (<div key={staple.id} className="archive-list">
          <div key={staple.id}>
            {/* eslint-disable-next-line jsx-a11y/anchor-is-valid */}
            <a onClick={() => fetchStaple(staple.id)} >{staple.name}</a><br/>
            {"Created: " + new Date(staple.created_at).toLocaleString()}
          </div>
        </div>)
      })
    }
    let s = undefined;
    if (staple.staple) {
      s = (
        <div className="Staples">
          <form>
            <FormGroup controlId="content" className="staple-view">
              <pre>{staple.staple.content}</pre>
            </FormGroup>
            <LoaderButton
              block
              bsSize="large"
              bsStyle="danger"
              onClick={() => handleDelete(staple.staple.id)}
              isLoading={isDeleting}
            >
              Delete
            </LoaderButton>
          </form>
          <br/>
        </div>
      )
    } else {
      s = (
        <div className="Staples">
          <form>
            <FormGroup controlId="content" className="staple-view">
              <pre>No archived staples.</pre>
            </FormGroup>
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
      </div>
    )
  }

  function renderArchives() {
    return (
      <div className="staples">
        <PageHeader>Your Archived Staples</PageHeader>
        <ListGroup>
          {!isLoading && renderArchiveView(staples)}
        </ListGroup>
      </div>
    );
  }

  return (
    <div className="Staples">
      {renderArchives()}
    </div>
  );
}