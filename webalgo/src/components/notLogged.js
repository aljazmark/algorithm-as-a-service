import React from 'react'
export const NotLogged = () => {
    return (
        <div className="flexDir" >
            <div className="divHelp" data-testit="requests" align="center">
                <h1 data-testid="requests-title">Please login to view this page</h1>
            </div>
        </div>
    )
}
