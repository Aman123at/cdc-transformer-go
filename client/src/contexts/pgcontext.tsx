'use client'

import { deleteRowApi, getTables, insertRow } from '@/apiCalls/commonCalls';
import { IRow, ITable } from '@/interfaces/commonInterface';
import React, { createContext, useContext, useEffect, useState } from 'react'
import {v4 as uuidv4} from "uuid"

export type PGContextType = {
  tables:  ITable[] | null;
  setTables: Function;
  fetchTables:Function;
  handleInsertRow:(tablename:string,newRow: IRow)=>void;
  handleDeleteRow:(tablename:string,rowid: number)=>void;
}

const PGContext = createContext<PGContextType | undefined>(undefined)

export function PGProvider({ children }: { children: React.ReactNode }) {
  const [tables, setTables] = useState<ITable[] | null>(null)
  const handleInsertRow = async (tablename:string,newRow: IRow) => {
    // call insert row api
    const {err} = await insertRow(tablename,newRow)
    if(err){
        console.log(err)
    }else{
        // call getTables api
        fetchTables()
    }
  }

  const handleDeleteRow = async (tablename:string,rowid:number) =>{
    const {err} = await deleteRowApi(tablename,rowid)
    if(err){
        console.log(err)
    }else{
        // call getTables api
        fetchTables()
    }
  }

  const fetchTables = async() =>{
    // get sessionid from localstorage
    let sessionId = localStorage.getItem("cdc-session-id")
    const {data,err}:any = await getTables(sessionId || "")
    if(err){
        setTables(null)
    }else{
        // console.log("TABLES>>",data)
        setTables(data)
    }
  }

 

  useEffect(()=>{
    fetchTables()
  },[])

  return (
    <PGContext.Provider value={{ tables, setTables, handleInsertRow, handleDeleteRow, fetchTables }}>
      {children}
    </PGContext.Provider>
  )
}

export function usePG() {
  const context = useContext(PGContext)
  if (context === undefined) {
    throw new Error('usePG must be used within a PGProvider')
  }
  return context
}