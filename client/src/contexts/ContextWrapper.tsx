import React from 'react'
import { PGProvider } from './pgcontext'
import { MongoProvider } from './mongocontext'

const ContextWrapper = ({children}:{children:React.ReactNode}) => {
  return (
    <PGProvider>
        <MongoProvider>
            {children}
        </MongoProvider>
    </PGProvider>
  )
}

export default ContextWrapper