import React from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import InboxIcon from '@material-ui/icons/MoveToInbox';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { makeStyles, useTheme } from '@material-ui/core/styles';
import LoginDialog from './login';
import AccountBoxIcon from '@material-ui/icons/AccountBox';
import {useSelector} from "react-redux";
import { NewRequest} from "./newRequest";
import { NewData} from "./newData";
import {Help} from "./help"
import {RequestsTable} from './requests';
import {DatasTable} from './datas';
import {Switch,Route,useHistory} from "react-router-dom"
import {NotLogged} from './notLogged'
import {Settings} from './settings'



const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  drawer: {
    [theme.breakpoints.up("md")]: {
      width: drawerWidth,
      flexShrink: 0,
    },
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
    [theme.breakpoints.up("md")]: {
      display: 'none',
    },
  },
  // necessary for content to be below app bar
  toolbar: theme.mixins.toolbar,
  drawerPaper: {
    width: drawerWidth,
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(3),
  },
  title: {
    flexGrow: 1,
  },
}));


function ResponsiveDrawer(props) {
  const { window } = props;
  const classes = useStyles();
  const theme = useTheme();
  const [ mobileOpen, setMobileOpen] = React.useState(false);
  const userID = useSelector((state)=> state.user.userID);
  const username = useSelector((state)=> state.user.username);
  const token = useSelector((state)=> state.user.token);
  const history = useHistory()
  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };
  var userLogged = (userID) ? true : false;
  const navList = [
    {
      text: "New Request",
      icon: <InboxIcon />,
      page: '/',
      show: true
    },
    {
      text: "New Data",
      icon: <InboxIcon />,
      page: 'newData',
      show: userLogged
    },
    {
      text: "Requests",
      icon: <InboxIcon />,
      page: '/requests',
      show: userLogged
    },
    {
      text: "Data",
      icon: <InboxIcon />,
      page: '/data',
      show: userLogged
    },
    {
      text: "Help",
      icon: <InboxIcon />,
      page: 'help',
      show: true
    }
  ];

 
  const drawer = (
    <div className="drawer" data-testid="drawer">
      <div className={classes.toolbar} />
      <Divider />
      <List>
          <ListItem >
            <ListItemIcon> <AccountBoxIcon /></ListItemIcon>
            <ListItemText>{username}</ListItemText>
          </ListItem>
          {userLogged ? <ListItem button onClick={()=>history.push('/settings')}>
          <ListItemIcon> <AccountBoxIcon /></ListItemIcon>
          <ListItemText>Settings</ListItemText>
        </ListItem> : null}
      </List>
      <Divider />
      <List>
        {navList.map((item, index) => {
          const { text, icon, page,show } = item;
          if (show) {
          return (
            <ListItem button key={text} onClick={()=>history.push(page)} >
              <ListItemIcon>{icon}</ListItemIcon>
              <ListItemText primary={text} />
            </ListItem>
          );
        }else{
          return null
        }
        })}
      </List>
    </div>
  );

  const container = window !== undefined ? () => window().document.body : undefined;

  return (
    <div className={classes.root}>
      <CssBaseline />
      <AppBar position="fixed" className={classes.appBar}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            className={classes.menuButton}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap className={classes.title}>
            Algo API
          </Typography>
          <LoginDialog />
        </Toolbar>
      </AppBar>
      <nav className={classes.drawer} aria-label="mailbox folders">
        {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
        <Hidden mdUp implementation="css">
          <Drawer
            container={container}
            variant="temporary"
            anchor={theme.direction === 'rtl' ? 'right' : 'left'}
            open={mobileOpen}
            onClose={handleDrawerToggle}
            classes={{
              paper: classes.drawerPaper,
            }}
            ModalProps={{
              keepMounted: true, // Better open performance on mobile.
            }}
          >
            {drawer}
          </Drawer>
        </Hidden>
        <Hidden smDown implementation="css">
          <Drawer
            classes={{
              paper: classes.drawerPaper,
            }}
            variant="permanent"
            open
          >
            {drawer}
          </Drawer>
        </Hidden>
      </nav>
      <main className={classes.content}>
        <div className={classes.toolbar} />
              <Switch>
                <Route exact path="/">
                  <NewRequest />
                </Route>
                <Route exact path="/help">
                  <Help />
                </Route>
                <Route exact path="/data">
                {userLogged ? <DatasTable userID={userID} token={token}/>: <NotLogged />}
                </Route>
                <Route exact path="/requests">
                {userLogged ? <RequestsTable userID={userID} token={token}/>: <NotLogged />}
                </Route>
                <Route exact path="/newData">
                {userLogged ? <NewData token={token}/>: <NotLogged />}
                </Route>
                <Route exact path="/settings">
                {userLogged ? <Settings userID={userID} token={token}/>: <NotLogged />}
                </Route>
              </Switch>
      </main>
    </div>
  );
}

ResponsiveDrawer.propTypes = {
  /**
   * Injected by the documentation to work in an iframe.
   * You won't need it on your project.
   *
   */
  window: PropTypes.func,
};

export default ResponsiveDrawer;