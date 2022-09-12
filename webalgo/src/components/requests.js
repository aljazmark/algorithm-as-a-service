import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import {Button} from '@material-ui/core';
import Collapse from '@material-ui/core/Collapse';
import IconButton from '@material-ui/core/IconButton';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp';
import {useDispatch,useSelector} from "react-redux";
import { getRequests,deleteRequest } from '../context';
const useRowStyles = makeStyles({
  root: {
    '& > *': {
      borderBottom: 'unset',
    },
  },
});

function Row(props) {
  const { row } = props;
  const [open, setOpen] = React.useState(false);
  const classes = useRowStyles();
  return (
    <React.Fragment>
      <TableRow className={classes.root}>
        <TableCell>
          <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {row.id}
        </TableCell>
        <TableCell >{row.algorithm}</TableCell>
        <TableCell >{row.requested}</TableCell>
        <TableCell align="right"> 
          <Button onClick={()=>{props.deleteRequest({id:row.id})}}>Delete</Button>
        </TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box margin={1}>
              <Typography variant="h6" gutterBottom component="div">
                Details
              </Typography>
              <Table size="small" aria-label="purchases">
                <TableHead>
                  <TableRow>
                    <TableCell overflow="auto" >Input</TableCell>
                    <TableCell >Parameters</TableCell>
                    <TableCell >Output</TableCell>
                    <TableCell >Execution Time</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                    <TableRow key={row.id}>
                      <TableCell component="th" scope="row">
                        {row.input}
                      </TableCell>
                      <TableCell>{JSON.stringify(row.parameters)}</TableCell>
                      <TableCell >{row.output}</TableCell>
                      <TableCell >
                        {row.executionTime}
                      </TableCell>
                    </TableRow>
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </React.Fragment>
  );
}

Row.propTypes = {
  row: PropTypes.shape({
    id: PropTypes.string.isRequired,
    algorithm: PropTypes.string.isRequired,
    requested: PropTypes.string.isRequired,
    input: PropTypes.string.isRequired,
    parameters: PropTypes.array,
    output: PropTypes.string.isRequired,
    executionTime:PropTypes.string.isRequired,
  }).isRequired,
};



export const RequestsTable = (props)=> {
  const requests = useSelector((state)=> state.user.requests);
  const dispatchR = useDispatch();
  const refreshRequests = () =>{
    getRequests(dispatchR,{userID:props.userID,token:props.token})
  } 
  const deleteRequestOption = (props2) =>{
    deleteRequest(dispatchR,{token:props.token,id:props2.id})
  } 
  return (  
    <div align="center">
      <h1 data-testid="requests-title">Your requests</h1>
      <TableContainer component={Paper}> 
        <Table aria-label="collapsible table">
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell>Request ID</TableCell>
              <TableCell>Algorithm</TableCell>
              <TableCell>Date</TableCell>
              <TableCell align="right" variant="head" size="small">
              <Button   onClick={refreshRequests}>
              Refresh
              </Button>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {requests.map((row) => (
              <Row key={row.id} row={row} deleteRequest={deleteRequestOption}/>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
}