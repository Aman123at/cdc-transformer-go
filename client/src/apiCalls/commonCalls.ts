import { IRow } from "@/interfaces/commonInterface";
import axios from "axios";
const API_BASE_URL = `${process.env.NEXT_PUBLIC_API_URL}`
type Table = {
    name: string
    columns: Column[]
    sessionid?:string
  }

  type Column = {
    name: string
    type: string
  }

  const apiWrapper = async <T>(
    apiCall: () => Promise<any>,
    errorMessage: string
): Promise<{ data: T | null; err: string | null }> => {
    try {
        const response = await apiCall();
        if (response.status === 200) {
            return { data: response.data, err: null };
        } else {
            return { data: null, err: errorMessage };
        }
    } catch (error) {
        return { data: null, err: `${errorMessage}: ${error}` };
    }
};

const createTable = async (tabledata: Table) => {
    return apiWrapper(
        () => axios.post(`${API_BASE_URL}/create/table`, { ...tabledata }),
        "Something went wrong while creating table"
    );
};

const getTables = async (sessionId:string) => {
    return apiWrapper(
        () => axios.get(`${API_BASE_URL}/fetch/tables/${sessionId}`),
        "Something went wrong while fetching tables"
    );
};

const getCollections = async (sessionId:string) => {
    return apiWrapper(
        () => axios.get(`${API_BASE_URL}/fetch/collections/${sessionId}`),
        "Something went wrong while fetching collections"
    );
};

const insertRow = async (tablename: string, rowData: IRow) => {
    return apiWrapper(
        () => axios.post(`${API_BASE_URL}/insert/row`, { tablename, row: rowData }),
        `Something went wrong while inserting row in table ${tablename}`
    );
};

const deleteRowApi = async (tablename: string, rowid: number) => {
    return apiWrapper(
        () => axios.post(`${API_BASE_URL}/delete/row`, { tablename, rowid }),
        `Something went wrong while deleting row ${rowid} in table ${tablename}`
    );
};



export {createTable,getTables,insertRow,deleteRowApi,getCollections};