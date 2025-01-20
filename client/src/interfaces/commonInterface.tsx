export interface IConstraint {
    default?:string;
    nullable:string;
    primary_key?:boolean;
}

export interface IColumn{
    name:string;
    type:string;
    constraint?:IConstraint
}

export interface IRow  {
    [key: string]: any
}

export interface ITable {
    tablename:string;
    rows:IRow[];
    columns:IColumn[];
}

export interface ICollection {
    collectionname:string;
    documents:{[key:string]:any}[];
}



