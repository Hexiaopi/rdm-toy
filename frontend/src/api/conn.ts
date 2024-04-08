import request from '@/utils/request'

export function listConns(query: any) {
    return request({
        url: '/api/v1/conns',
        method: 'get',
        params: query
    })
}

export function getConn(conn: string, query: any) {
    return request({
        url: `/api/v1/conn/${conn}`,
        method: 'get',
        params: query
    })
}

export function createConn(data: any) {
    return request({
        url: '/api/v1/conn',
        method: 'post',
        data: data,
    })
}

export function deleteConn(conn: string) {
    return request({
        url: `/api/v1/conn/${conn}`,
        method: 'delete',
    })
}

export function testConn(data: any) {
    return request({
        url: '/api/v1/conn/test',
        method: 'post',
        data: data,
    })
}

export function listClients(conn: string) {
    return request({
        url: `/api/v1/conn/${conn}/clients`,
        method: 'get',
    })
}

export function deleteClient(conn: string, data: any) {
    return request({
        url: `/api/v1/conn/${conn}/client`,
        method: 'delete',
        data: data,
    })
}