package animation

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Dimension displays the dimension change screen to the player.
type Dimension struct{}

// Play ...
func (animation *Dimension) Play(conn *minecraft.Conn, serverGameData minecraft.GameData) {
	var dimension int32
	if conn.GameData().Dimension == packet.DimensionNether {
		dimension = packet.DimensionEnd
	} else {
		dimension = packet.DimensionNether
	}
	sendDimension(conn, serverGameData, dimension, false)
}

// Clear stops the dimension change animation and safely updates the dimension to the destination.
func (animation *Dimension) Clear(conn *minecraft.Conn, serverGameData minecraft.GameData) {
	sendDimension(conn, serverGameData, serverGameData.Dimension, true)
}

// sendDimension updates the player's dimension and sends the client-side dimension switch ack
// using the player's actual entity runtime ID, followed by player spawn if requested.
func sendDimension(conn *minecraft.Conn, serverGameData minecraft.GameData, dimension int32, playStatus bool) {
	_ = conn.WritePacket(&packet.ChangeDimension{Dimension: dimension, Position: serverGameData.PlayerPosition})
	_ = conn.WritePacket(&packet.StopSound{StopAll: true})
	_ = conn.WritePacket(&packet.PlayerAction{
		EntityRuntimeID: serverGameData.EntityRuntimeID,
		ActionType:      protocol.PlayerActionDimensionChangeDone,
	})
	if playStatus {
		_ = conn.WritePacket(&packet.PlayStatus{
			Status: packet.PlayStatusPlayerSpawn,
		})
	}
}
