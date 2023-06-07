// Copyright 2023 Blink Labs, LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chainsync

import (
	"github.com/blinklabs-io/snek/plugin"
)

var cmdlineOptions struct {
	network        string
	networkMagic   uint
	address        string
	socketPath     string
	ntcTcp         bool
	intersectTip   bool
	intersectPoint string
	includeCbor    bool
}

func init() {
	plugin.Register(
		plugin.PluginEntry{
			Type:               plugin.PluginTypeInput,
			Name:               "chainsync",
			Description:        "syncs blocks from a Cardano node using either NtC (node-to-client) or NtN (node-to-node)",
			NewFromOptionsFunc: NewFromCmdlineOptions,
			Options: []plugin.PluginOption{
				{
					Name:         "network",
					Type:         plugin.PluginOptionTypeString,
					CustomEnvVar: "CARDANO_NETWORK",
					Description:  "specifies a well-known Cardano network name",
					DefaultValue: "mainnet",
					Dest:         &(cmdlineOptions.network),
				},
				{
					Name:         "network-magic",
					Type:         plugin.PluginOptionTypeUint,
					Description:  "specifies the network magic value to use, overrides 'network'",
					DefaultValue: uint(0),
					Dest:         &(cmdlineOptions.networkMagic),
				},
				{
					Name:         "address",
					Type:         plugin.PluginOptionTypeString,
					Description:  "specifies the TCP address of the node to connect to in the form 'host:port'",
					DefaultValue: "",
					Dest:         &(cmdlineOptions.address),
				},
				{
					Name:         "socket-path",
					Type:         plugin.PluginOptionTypeString,
					CustomEnvVar: "CARDANO_NODE_SOCKET_PATH",
					Description:  "specifies the path to the UNIX socket to connect to",
					DefaultValue: "",
					Dest:         &(cmdlineOptions.socketPath),
				},
				{
					Name:         "ntc-tcp",
					Type:         plugin.PluginOptionTypeBool,
					Description:  "use the NtC (node-to-client) protocol over TCP, for use when exposing a node's UNIX socket via socat or similar",
					DefaultValue: false,
					Dest:         &(cmdlineOptions.ntcTcp),
				},
				{
					Name:         "intersect-tip",
					Type:         plugin.PluginOptionTypeBool,
					Description:  "start syncing at the chain tip (defaults to chain genesis)",
					DefaultValue: true,
					Dest:         &(cmdlineOptions.intersectTip),
				},
				// TODO: intersect-point
				{
					Name:         "include-cbor",
					Type:         plugin.PluginOptionTypeBool,
					Description:  "include original CBOR for block/transaction in events",
					DefaultValue: false,
					Dest:         &(cmdlineOptions.includeCbor),
				},
			},
		},
	)
}

func NewFromCmdlineOptions() plugin.Plugin {
	p := New(
		WithNetwork(cmdlineOptions.network),
		WithNetworkMagic(uint32(cmdlineOptions.networkMagic)),
		WithAddress(cmdlineOptions.address),
		WithSocketPath(cmdlineOptions.socketPath),
		WithNtcTcp(cmdlineOptions.ntcTcp),
		WithIntersectTip(cmdlineOptions.intersectTip),
		// TODO: WithIntersectPoints
		WithIncludeCbor(cmdlineOptions.includeCbor),
	)
	return p
}
