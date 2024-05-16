from itertools import count
from pathlib import Path
from re import search
from typing import List
from xml.etree.ElementTree import Element, ElementTree

import click
import pandas
from lxml import etree
from matplotlib.colors import XKCD_COLORS
from onnx import load
from onnx.onnx_pb import GraphProto, ModelProto, NodeProto
from pandas import DataFrame
from progress.bar import Bar

NODE_ID_COUNTER: count = count()
EDGE_ID_COUNTER: count = count()
OUTPUT_DF_LIST: List[DataFrame] = []


def extractLayer(nodeName: str) -> str:
    """


    :param nodeName: str:

    """
    pattern: str = r"layer\.(\d+)"

    layer: str
    try:
        layer = search(pattern=pattern, string=nodeName).group(0)
    except AttributeError:
        layer = ""

    return layer


def buildDF(
    nodeID: int,
    name: str,
    opType: str,
    layer: str,
    inputs: List[str],
    outputs: List[str],
    color: List[str],
) -> DataFrame:
    """


    :param nodeID: int:
    :param name: str:
    :param opType: str:
    :param layer: str:
    :param inputs: List[str]:
    :param outputs: List[str]:
    :param color: List[str]:

    """
    data: dict[str, List[int | str | List[str]]] = {
        "ID": [nodeID],
        "Name": [name],
        "Op Type": [opType],
        "Layer": [layer],
        "Inputs": [inputs],
        "Outputs": [outputs],
        "Color": [color],
    }
    return DataFrame(data)


def dfIDQuery(df: DataFrame, query: str) -> tuple[str, str] | None:
    """


    :param df: DataFrame:
    :param query: str:

    """
    mask = df["Outputs"].apply(lambda x: query in x)
    tempDF: DataFrame = df[mask]

    try:
        return (tempDF["Name"].iloc[0], str(tempDF["ID"].iloc[0]))
    except IndexError:
        return None


def buildXML(
    df: DataFrame,
    mode: str = "production",
) -> str:
    """


    :param df: DataFrame:
    :param mode: str:  (Default value = "production")

    """
    edgeList: List[tuple[tuple[str, str], str]] = []

    version: str
    if mode == "production":
        version = "1.2draft"
    else:
        version = "1.2"

    xmlns: str = f"http://www.gexf.net/{version}"
    xmlnsViz: str = f"http://gexf.net/{version}/viz"

    rootNode: Element = etree.Element("gexf", nsmap={None: xmlns, "viz": xmlnsViz})
    rootNode.set("version", version)

    graphNode: Element = etree.SubElement(rootNode, "graph")
    graphNode.set("mode", "static")
    graphNode.set("defaultedgetype", "directed")
    graphNode.set("idtype", "integer")

    attributesNode: Element = etree.SubElement(graphNode, "attributes")
    attributesNode.set("class", "node")

    inputAttributeNode: Element = etree.SubElement(attributesNode, "attribute")
    inputAttributeNode.set("id", "type")
    inputAttributeNode.set("title", "Operation_Type")
    inputAttributeNode.set("type", "string")

    inputAttributeNode: Element = etree.SubElement(attributesNode, "attribute")
    inputAttributeNode.set("id", "input")
    inputAttributeNode.set("title", "Input")
    inputAttributeNode.set("type", "string")

    outputAttributeNode: Element = etree.SubElement(attributesNode, "attribute")
    outputAttributeNode.set("id", "output")
    outputAttributeNode.set("title", "Output")
    outputAttributeNode.set("type", "string")

    layerAttributeNode: Element = etree.SubElement(attributesNode, "attribute")
    layerAttributeNode.set("id", "layer")
    layerAttributeNode.set("title", "Layer")
    layerAttributeNode.set("type", "string")

    verticesNode: Element = etree.SubElement(graphNode, "nodes")
    edgesNode: Element = etree.SubElement(graphNode, "edges")

    with Bar("Creating GEXF nodes...", max=df.shape[0]) as bar:
        for ID, NAME, OPTYPE, LAYER, INPUTS, OUTPUTS, COLOR in df.itertuples(
            index=False
        ):
            ID: str = str(ID)
            vertexNode: Element = etree.SubElement(verticesNode, "node")
            vertexNode.set("id", ID)
            vertexNode.set("label", NAME)

            vizColorNode: Element = etree.SubElement(vertexNode, "color")
            vizColorNode.set("hex", COLOR)

            attvaluesNode: Element = etree.SubElement(vertexNode, "attvalues")

            attvalueNode: Element = etree.SubElement(attvaluesNode, "attvalue")
            attvalueNode.set("for", "type")
            attvalueNode.set("value", OPTYPE)

            attvalueNode: Element = etree.SubElement(attvaluesNode, "attvalue")
            attvalueNode.set("for", "layer")
            attvalueNode.set("value", LAYER)

            i: str
            for i in INPUTS:
                attvalueNode: Element = etree.SubElement(attvaluesNode, "attvalue")
                attvalueNode.set("for", "input")
                attvalueNode.set("value", i.replace(":", "-"))

                parentNodeNameID: tuple[str, str] | None = dfIDQuery(df=df, query=i)

                if parentNodeNameID is None:
                    pass
                else:
                    nodePairing: tuple[tuple[str, str], str] = (parentNodeNameID, ID)
                    edgeList.append(nodePairing)

            o: str
            for o in OUTPUTS:
                attvalueNode: Element = etree.SubElement(attvaluesNode, "attvalue")
                attvalueNode.set("for", "output")
                attvalueNode.set("value", o.replace(":", "-"))

            bar.next()

    with Bar("Creating GEXF edges...", max=len(edgeList)) as bar:
        pair: tuple[tuple[str, str], str]
        for pair in edgeList:
            edgeNode: Element = etree.SubElement(edgesNode, "edge")
            edgeNode.set("id", str(EDGE_ID_COUNTER.__next__()))
            edgeNode.set("source", pair[0][1])
            edgeNode.set("target", pair[1])
            edgeNode.set("label", pair[0][0])
            bar.next()

    xmlStr: str = etree.tostring(rootNode, pretty_print=True).decode()
    xmlStr = xmlStr.replace("<color", "<viz:color")

    return xmlStr


@click.command()
@click.option(
    "onnxFile",
    "-i",
    "--input",
    nargs=1,
    type=Path,
    required=True,
    help="Path to ONNX file",
)
@click.option(
    "gexfFile",
    "-o",
    "--output",
    nargs=1,
    type=Path,
    required=True,
    help="Path to store GEXF file",
)
@click.option(
    "mode",
    "-m",
    "--mode",
    nargs=1,
    type=click.Choice(choices=["production", "validation"]),
    default="production",
    show_default=True,
    help="Save the GEXF for production usage or for validation usage",
)
def main(onnxFile: Path, gexfFile: Path, mode: str) -> None:
    """
    Convert an ONNX computational graph into a GEXF formatted file
    \f

    :param onnxFile: Path:
    :param gexfFile: Path:
    :param mode: str:

    """
    colors: List[str] = list(XKCD_COLORS.values())

    model: ModelProto = load(f=onnxFile.__str__())
    graph: GraphProto = model.graph

    with Bar(
        "Extracting information from ONNX computational graph...", max=len(graph.node)
    ) as bar:
        previousLayer: str = ""
        colorIDX: int = 0

        node: NodeProto
        for node in graph.node:
            nodeID: int = NODE_ID_COUNTER.__next__()
            name: str = node.name
            opType: str = node.op_type
            layer: str = extractLayer(nodeName=name)

            if layer != previousLayer:
                colorIDX += 1
                previousLayer = layer

            color: str = colors[colorIDX]
            outputs: List[str] = list(node.output)
            inputs: List[str] = list(node.input)
            df: DataFrame = buildDF(
                nodeID=nodeID,
                name=name,
                opType=opType,
                layer=layer,
                inputs=inputs,
                outputs=outputs,
                color=color,
            )
            OUTPUT_DF_LIST.append(df)
            bar.next()
    df: DataFrame = pandas.concat(OUTPUT_DF_LIST)
    xmlStr = buildXML(df=df, mode=mode)

    with open(file=gexfFile, mode="w") as xmlFile:
        xmlFile.write(xmlStr)
        xmlFile.close()


if __name__ == "__main__":
    main()
