import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import Box from '@material-ui/core/Box';
import Collapse from '@material-ui/core/Collapse';
import {Button} from '@material-ui/core';
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
import { getDatas,deleteData } from '../context';
import {DataDialog} from './dataDialog'
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
        <TableCell >{row.user}</TableCell>
        <TableCell >{row.created}</TableCell>
        <TableCell align="right"> 
          <DataDialog dataID={row.id} data={row.data} allowed={row.allowed}/>
          <Button onClick={()=>{props.deleteData({id:row.id})}}>Delete</Button>
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
                    <TableCell overflow="auto" >Data</TableCell>
                    <TableCell >Allowed</TableCell>
                    <TableCell >Updated</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                    <TableRow key={row.id}>
                      <TableCell component="th" scope="row">
                        {row.data}
                      </TableCell>
                      <TableCell>{JSON.stringify(row.allowed)}</TableCell>
                      <TableCell >{row.updated}</TableCell>
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
    user: PropTypes.string.isRequired,
    created: PropTypes.string.isRequired,
    data: PropTypes.string.isRequired,
    allowed: PropTypes.array,
    updated: PropTypes.string.isRequired,
  }).isRequired,
};



export const DatasTable = (props)=> {
  const datas = useSelector((state)=> state.user.datas);
  const dispatchR = useDispatch();
  const refreshData = () =>{
    getDatas(dispatchR,{userID:props.userID,token:props.token})
  } 
  const deleteDataOption = (props2) =>{
    deleteData(dispatchR,{token:props.token,id:props2.id})
  } 
  return (        
    <div align="center">
      <h1 data-testid="requests-title">Your data</h1>
      <TableContainer component={Paper}>
        <Table aria-label="collapsible table">
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell>Data ID</TableCell>
              <TableCell>User</TableCell>
              <TableCell>Date</TableCell>
              <TableCell align="right" variant="head" size="small">
              <Button   onClick={refreshData}>
              Refresh
              </Button>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {datas.map((row) => (
              <Row key={row.id} row={row} deleteData={deleteDataOption}/>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
}