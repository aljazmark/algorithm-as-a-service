import React from 'react'
import {useState,useEffect} from 'react'
import { TextField,MenuItem,Select,InputLabel,makeStyles,FormControl } from '@material-ui/core';
import {useDispatch,useSelector} from "react-redux";
import store from '../context/store'
import { getAlgorithmHelp } from '../context';
import {isEmpty} from "lodash"
const useStyles = makeStyles((theme) => ({
    formControl: {
      margin: theme.spacing(0),
      minWidth: 210,
      minHeight: 150,
    },
    selectEmpty: {
      marginTop: theme.spacing(0),
    },
    resultText:{
        minWidth:300
    }
  }));

 


const initialFormValues ={
    algorithm:""
}

export const Help = () => {
    const classes = useStyles();
    const [formValues,setFormValues] = useState(initialFormValues);
    const [details,setDetails] = useState({});
    const [show,setShow] = useState(false);
    const algorithms = useSelector((state)=> state.user.algorithms);
    const dispatchR = useDispatch();
    const handleHelp= e =>{
        const { name,value} = e.target
        setFormValues({
            ...formValues,
            [name]:value
        })
    }
    useEffect(() => {
        var found = false
        if(formValues.algorithm!==""){
            var helps = store.getState().user.helps
            helps.forEach(function (arrayItem){
                if(formValues.algorithm === arrayItem.algorithm){
                    setDetails(arrayItem)
                    found= true;
                    return;
                }
            })
            if(!found){
                getAlgorithmHelp(dispatchR,{algorithm:formValues.algorithm}).then(res=> {setDetails(res)})
            }
            
        }
    }, [formValues]); // eslint-disable-line react-hooks/exhaustive-deps
    useEffect(() => {
        if(!isEmpty(details)){
            setShow(true)
        }
    }, [details]);
    const showDetails = () => {
        return (
        <div>
            <h3>Algorithm</h3>
                <TextField
                className={classes.resultText}
                variant="outlined"
                value={details.algorithm}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Category</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                value={details.category}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Description</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.description}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Input format</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.inputFormat}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Input Example</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.inputExample}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Parameters format</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.parametersFormat}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Parameters example</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.parametersExample}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Output format</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.outputFormat}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
            <h3>Output example</h3>
            <TextField
                className={classes.resultText}
                variant="outlined"
                multiline
                value={details.outputExample}
                InputProps={{
                    readOnly: true,
                }}
            ></TextField>
        </div>
                )
    }
    return (
        <div className="flexDir" >
            <div className="divHelp" data-testit="requests" align="center">
                <h1 data-testid="requests-title">Help</h1>
                <form>
                    <FormControl variant="outlined">
                        <InputLabel>Select an algorithm</InputLabel>
                        <Select
                        variant="outlined"
                        label="Select an algorithm"
                        name="algorithm"
                        className={classes.resultText}
                        value={formValues.algorithm}
                        onChange= {handleHelp}>
                            {
                                algorithms.map(
                                    item => (<MenuItem key={item} value={item}>{item}</MenuItem>)

                                )
                            }
                            </Select>
                    </FormControl>    
                </form>
                <div>
                    <h2 data-testid="requests-title">Algorithm details</h2>
                </div>
                {show ? showDetails(): null}
            </div>
        </div>
    )
}
