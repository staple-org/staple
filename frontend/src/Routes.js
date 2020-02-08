import React from "react";
import { Route, Switch } from "react-router-dom";
import Home from "./containers/Home";
import Login from "./containers/Login";
import NotFound from "./containers/NotFound";
import AppliedRoute from "./components/AppliedRoute";
import Signup from "./containers/Signup";
import AddStaple from "./containers/AddStaple";
import Settings from "./containers/Settings";
import AuthenticatedRoute from "./components/AuthenticatedRoute";
import UnauthenticatedRoute from "./components/UnauthenticatedRoute";
import Archives from "./containers/Archives";
import ResetPassword from "./containers/ResetPassword";


export default function Routes({ appProps }) {
  return (
    <Switch>
      <AppliedRoute path="/" exact component={Home} appProps={appProps} />
      <UnauthenticatedRoute path="/login" exact component={Login} appProps={appProps} />
      <UnauthenticatedRoute path="/signup" exact component={Signup} appProps={appProps} />
      <UnauthenticatedRoute path="/reset" exact component={ResetPassword} appProps={appProps} />
      <AuthenticatedRoute path="/settings" exact component={Settings} appProps={appProps} />
      <AuthenticatedRoute path="/staples/new" exact component={AddStaple} appProps={appProps} />
      <AuthenticatedRoute path="/archive" exact component={Archives} appProps={appProps} />
      <Route component={NotFound} />
    </Switch>
  );
}