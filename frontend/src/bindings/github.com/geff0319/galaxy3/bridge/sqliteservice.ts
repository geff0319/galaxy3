// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as sql$0 from "../../../../database/sql/models.js";

export function Close(): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3096589742) as any;
    return $resultPromise;
}

export function Execute(query: string, ...args: any[]): Promise<sql$0.Result> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3791336209, query, args) as any;
    return $resultPromise;
}

export function Init(): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(4025795672) as any;
    return $resultPromise;
}

/**
 * Name returns the name of the plugin.
 * You should use the go module format e.g. github.com/myuser/myplugin
 */
export function Name(): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(123686969) as any;
    return $resultPromise;
}

export function Open(dbPath: string): Promise<string> & { cancel(): void } {
    let $resultPromise = $Call.ByID(221954298, dbPath) as any;
    return $resultPromise;
}

export function Select(query: string, ...args: any[]): Promise<({ [_: string]: any } | null)[] | null> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1953141994, query, args) as any;
    return $resultPromise;
}
