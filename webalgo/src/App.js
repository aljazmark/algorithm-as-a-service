import React from 'react';
import ResponsiveDrawer from './components/navigation'
import { createMuiTheme, ThemeProvider } from '@material-ui/core';
import { getDatas,getAlgorithms,getRequests, getUser } from './context';
import {useDispatch,useSelector} from "react-redux";
import {BrowserRouter as Router} from "react-router-dom"
const theme = createMuiTheme({
  palette: {
      primary: {
        main: '#b71c1c',
      },
      secondary: {
        main: '#b71c1c',
      },
    },
});
export const App = () => {
  const userID = useSelector((state)=> state.user.userID);
  const token = useSelector((state)=> state.user.token);
  const dispatchR = useDispatch();
    getAlgorithms(dispatchR)
  if(userID){
    getRequests(dispatchR,{userID:userID,token:token})
    getDatas(dispatchR,{userID:userID,token:token})
    getUser(dispatchR,{userID:userID,token:token})
  }
  return(
    <Router>
      <ThemeProvider theme={theme}>
          <ResponsiveDrawer />
      </ThemeProvider>
    </Router>
  )
  };
export default App;
