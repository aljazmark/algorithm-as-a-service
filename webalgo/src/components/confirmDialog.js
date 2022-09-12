import React from 'react'
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, Typography } from '@material-ui/core'
export const ConfirmDialog = (props) => {
    return (
        <Dialog open={props.confirm}>
            <DialogTitle>
                {props.title}
            </DialogTitle>
            <DialogContent>
                <Typography>
                    {props.content}
                </Typography>
            
            </DialogContent>
            <DialogActions>
                <Button 
                color="default"
                onClick={()=>props.setConfirm(false)}
                >
                Cancel
                </Button>

                <Button 
                color="primary"
                onClick={()=>props.confirmChange(true)}
                >
                    Confirm
                </Button>
            </DialogActions>
        </Dialog>
    )
}
