'use client'

import { getCollections } from '@/apiCalls/commonCalls';
import { ICollection } from '@/interfaces/commonInterface';

import React, { createContext, useContext, useEffect, useState } from 'react'


export type MongoContextType = {
  collections:  ICollection[] | null;
  setCollections: Function;
}

const MongoContext = createContext<MongoContextType | undefined>(undefined)

export function MongoProvider({ children }: { children: React.ReactNode }) {
  const [collections, setCollections] = useState<ICollection[] | null>(null)

  const fetchCollections = async() =>{
    let sessionid = localStorage.getItem("cdc-session-id")
    const {data,err}:any = await getCollections(sessionid || "")
    if(err){
        setCollections(null)
    }else{
        // console.log("COLLECTIONS>>",data)
        setCollections(data)
    }
  }

 

  useEffect(()=>{
    fetchCollections()
    const interval = setInterval(()=>fetchCollections(),10000)
    return ()=>clearInterval(interval)
  },[])

  return (
    <MongoContext.Provider value={{ collections, setCollections }}>
      {children}
    </MongoContext.Provider>
  )
}

export function useMongo() {
  const context = useContext(MongoContext)
  if (context === undefined) {
    throw new Error('useMongo must be used within a MongoProvider')
  }
  return context
}