/*
* Licensed to the Apache Software Foundation (ASF) under one or more
* contributor license agreements.  See the NOTICE file distributed with
* this work for additional information regarding copyright ownership.
* The ASF licenses this file to You under the Apache License, Version 2.0
* (the "License"); you may not use this file except in compliance with
* the License.  You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package xr

//
//import (
//	"context"
//
//	"mosn.io/pkg/log"
//
//	buf "mosn.io/mosn/pkg/buffer"
//	"mosn.io/pkg/buffer"
//)
//
//var bufCtx xrBufferCtx
//
//func init() {
//	buf.RegisterBuffer(&bufCtx)
//}
//
//type xrBufferCtx struct {
//	buf.TempBufferCtx
//}
//
//func (ctx xrBufferCtx) New() interface{} {
//	return new(xrBuffer)
//}
//
//func (ctx xrBufferCtx) Reset(i interface{}) {
//	buf, _ := i.(*xrBuffer)
//
//	// recycle ioBuffer
//	if buf.request.Data != nil {
//		if e := buffer.PutIoBuffer(buf.request.Data); e != nil {
//			log.DefaultLogger.Errorf("[protocol] [xr] [buffer] [reset] PutIoBuffer error: %v", e)
//		}
//	}
//
//	if buf.response.Data != nil {
//		if e := buffer.PutIoBuffer(buf.response.Data); e != nil {
//			log.DefaultLogger.Errorf("[protocol] [xr] [buffer] [reset] PutIoBuffer error: %v", e)
//		}
//	}
//
//	*buf = xrBuffer{}
//}
//
//type xrBuffer struct {
//	request  Request
//	response Response
//}
//
//func bufferByContext(ctx context.Context) *xrBuffer {
//	poolCtx := buf.PoolContext(ctx)
//	return poolCtx.Find(&bufCtx, nil).(*xrBuffer)
//}
