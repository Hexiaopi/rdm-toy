import request from '@/utils/request'

export function listKeys(conn: string, db: string, query: any) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/keys`,
        method: 'get',
        params: query
    })
}

export function getKey(conn: string, db: string, key: string) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/key/${key}`,
        method: 'get',
    })
}

export function addKey(conn: string, db: string, data: any) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/key`,
        method: 'post',
        data: data,
    })
}

export function deleteKey(conn: string, db: string, key: string) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/key/${key}`,
        method: 'delete',
    })
}

export function deleteKeys(conn: string, db: string) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/keys`,
        method: 'delete',
    })
}

export function deleteKeyField(conn: string, db: string, key: string, data: any) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/key/${key}/field`,
        method: 'delete',
        data: data,
    })
}


export function patchKeyName(conn: string, db: string, key: string, name: string) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/key/${key}/name`,
        method: 'patch',
        data: { "name": name },
    })
}

export function patchKeyTTL(conn: string, db: string, key: string, ttl: number) {
    return request({
        url: `/api/v1/conn/${conn}/db/${db}/key/${key}/ttl`,
        method: 'patch',
        data: { "ttl": ttl },
    })
}