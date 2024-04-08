import request from '@/utils/request'

export function listDBs(id: string, query: any) {
    return request({
        url: `/api/v1/conn/${id}/dbs`,
        method: 'get',
        params: query
    })
}